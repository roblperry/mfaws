package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/roblperry/mfaws/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "mfaws",
	Short: "mfaws - Multi-factors Authentication - AWS",
	Long:  `A little tool to help with updating session profiles generated using aws sts get-session-token`,
	Run:   rootCmdFunction,
}

/**
 * Init cobra (commander)
 */
func init() {
	cobra.OnInitialize(config.InitConfig)

	flagSet := rootCmd.PersistentFlags()
	flagSet.StringP("logging-level", "l", "Info", "Set the Logging Level to (Panic, Fatal, Error, Warn, Info, Debug, Trace)")
	err := viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("logging-level"))
	if err != nil {
		panic(fmt.Errorf("Fatal error binding logging level: %s \n", err))
	}

	flagSet.StringP("profile", "p", "", "Set the aws profile")
	err = viper.BindPFlag("aws.profile", rootCmd.PersistentFlags().Lookup("profile"))
	if err != nil {
		panic(fmt.Errorf("Fatal error binding aws profile: %s \n", err))
	}

	flagSet.StringP("region", "r", "", "Set the aws region")
	err = viper.BindPFlag("aws.region", rootCmd.PersistentFlags().Lookup("region"))
	if err != nil {
		panic(fmt.Errorf("Fatal error binding aws region: %s \n", err))
	}

	flagSet.StringP("session-profile", "s", "", "Session profile to configure with the acquired credentials")
	err = viper.BindPFlag("aws.session.profile", rootCmd.PersistentFlags().Lookup("session-profile"))
	if err != nil {
		panic(fmt.Errorf("Fatal error binding aws session profile: %s \n", err))
	}

	flagSet.StringVarP(&mfa, "mfa", "m", "", "MFA to use when obtaining session credentials")
}

//noinspection GoUnusedParameter
func rootCmdFunction(cmd *cobra.Command, args []string) {
	log.Info("Starting")

	if len(config.SessionProfile) == 0 {
		_, _ = fmt.Fprintf(os.Stderr,
			"I do not know the name of the profile to update or create.\n"+
				"Try setting profile or session-profile\n")
		os.Exit(1)
	}

	identity := getIdentity()

	arn := identity.Arn
	sn := strings.Replace(*arn, ":user/", ":mfa/", 1)

	mfa := getMFA()
	log.Debugf("%s is the MFA\n", mfa)

	token := getSessionToken(sn, mfa)

	if len(config.SessionProfile) == 0 {
		panic("Session Profile should be set")
	}

	setAccessKeyId(config.SessionProfile, *token.Credentials.AccessKeyId)
	setSecretAccessKey(config.SessionProfile, *token.Credentials.SecretAccessKey)
	setSessionToken(config.SessionProfile, *token.Credentials.SessionToken)

	_, _ = fmt.Fprintf(os.Stdout, "%s profile updated\n", config.SessionProfile)
	log.Info("Done")
}

func setSessionToken(sessionProfile string, sessionToken string) {
	cmd := exec.Command("aws", "configure", "--profile", sessionProfile, "set", "aws_session_token", sessionToken)
	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to set session token: %s\n", err)
	}
}

func setSecretAccessKey(sessionProfile string, secretAccessKey string) {
	cmd := exec.Command("aws", "configure", "--profile", sessionProfile, "set", "aws_secret_access_key", secretAccessKey)
	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to set session token: %s\n", err)
	}
}

func setAccessKeyId(sessionProfile string, accessKeyId string) {
	cmd := exec.Command("aws", "configure", "--profile", sessionProfile, "set", "aws_access_key_id", accessKeyId)
	err := cmd.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to set session token: %s\n", err)
	}
}

func getSessionToken(sn string, mfa string) *sts.GetSessionTokenOutput {
	input := &sts.GetSessionTokenInput{
		DurationSeconds: aws.Int64(129600),
		SerialNumber:    aws.String(sn),
		TokenCode:       aws.String(mfa),
	}
	token, err := getSts().GetSessionToken(input)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Authenticiation failed : %s\n", err)
		os.Exit(1)
	}
	return token
}

var mfa string

func getMFA() string {
	if len(mfa) == 0 {
		fmt.Print("Enter MFA: ")
		n, err := fmt.Scanf("%s", &mfa)
		if n != 1 {
			_, _ = fmt.Fprintf(os.Stderr, "scanned an unpexpected number of values %d\n", n)
			os.Exit(1)
		}

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed scanning mfa : %s\n", err)
			os.Exit(1)
		}
	}

	return mfa
}

var sess *session.Session

func getSession() *session.Session {
	if sess == nil {
		var err error
		sess, err = session.NewSession()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to get aws session: %s \n", err)
			os.Exit(1)
		}
	}

	return sess
}

var stsInstance *sts.STS

func getSts() *sts.STS {
	if stsInstance == nil {
		stsInstance = sts.New(getSession())
	}

	return stsInstance
}

var identity *sts.GetCallerIdentityOutput

func getIdentity() *sts.GetCallerIdentityOutput {
	var err error
	identity, err = getSts().GetCallerIdentity(nil)
	if err != nil {
		if e, ok := err.(awserr.Error); ok {
			if e.Code() == "NoCredentialProviders" {
				_, _ = fmt.Fprintln(os.Stderr, "I don't know how to authenticate")
				os.Exit(1)
			}
		}

		_, _ = fmt.Fprintf(os.Stderr, "Failed to get caller identity: %s \n", err)
		os.Exit(1)
	}

	return identity
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
