package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/roblperry/mfaws/cmd/values"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

/**
 * Init viper (configuration)
 */
func init() {
	var appName = AppName()
	var configFileFound = true

	viper.SetConfigName(appName)                     // name of config file (without extension)
	viper.AddConfigPath(path.Join("/etc", appName))  // path to look for the config file in
	viper.AddConfigPath(path.Join("$HOME", appName)) // call multiple times to add many search paths
	viper.AddConfigPath(".")                         // optionally look for config in the working directory
	err := viper.ReadInConfig()                      // Find and read the config file
	if err != nil {                                  // Handle errors reading the config file
		configFileFound = false
	}

	viper.SetEnvPrefix(appName)
	viper.AutomaticEnv()

	// Use - in config names, but _ in env names
	replacer := strings.NewReplacer("-", "_", ".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("logging.level", "ERROR")

	if configFileFound {
		viper.WatchConfig()
		viper.OnConfigChange(configChange)
	}
}

func InitConfig() {
	processConfig()
}

func configChange(_ fsnotify.Event) {
	processConfig()
}

func processConfig() {
	professLoggingLevel()
	processProfile()
	processRegion()
	processTargetProfile()
}

func logSetEnvError(key string, err error) {
	_, _ = fmt.Fprintf(os.Stderr, "Failed to set %s: %s\n", key, err)
	os.Exit(1)
}

func processProfile() {
	profile := viper.GetString("aws.profile")
	if len(profile) > 0 {
		if err := os.Setenv("AWS_PROFILE", profile); err != nil {
			logSetEnvError("AWS_PROFILE", err)
		}
		log.Infof("AWS_PROFILE now set to %s", profile)
	}
}

func processRegion() {
	region := viper.GetString("aws.region")
	if len(region) > 0 {
		if err := os.Setenv("AWS_REGION", region); err != nil {
			logSetEnvError("AWS_REGION", err)
		}
		log.Infof("AWS_REGION now set to %s", region)
	}
}

var TargetProfile string

func processTargetProfile() {
	TargetProfile = viper.GetString("aws.target_profile")
	if len(TargetProfile) == 0 {
		profile := os.Getenv("AWS_PROFILE")
		if len(profile) > 0 {
			TargetProfile = profile + "_session"
		}
	}

	if len(TargetProfile) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Neither profile nor target profile set\n")
		os.Exit(1)
	}

	log.Infof("Targeting profile %s", TargetProfile)
}

func professLoggingLevel() {
	loggingLevel := values.LoggingLevelValue{}
	err := loggingLevel.Set(viper.GetString("logging.level"))
	if err != nil {
		panic(fmt.Errorf("Fatal error reading logging level: %s \n", err))
	}
	log.SetLevel(loggingLevel.Level)
	log.Infof("Logging Level now set to %v", loggingLevel.Level)
}
