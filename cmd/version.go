package cmd

import (
	"fmt"
	"github.com/roblperry/mfaws/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of mfaws",
	Long:  `All software has versions. This is mfaws's`,
	Run:   versionCmdFunc,
}

//noinspection GoUnusedParameter
func versionCmdFunc(cmd *cobra.Command, args []string) {
	fmt.Println(config.AppName(), config.Version())
}
