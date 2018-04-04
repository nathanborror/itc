package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of itc",
	Long:  "All software has versions. This is itc's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(versionString())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionString() string {
	return fmt.Sprintf("itc v%s\n", version)
}
