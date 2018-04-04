package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testflightCmd = &cobra.Command{
	Use:   "testflight",
	Short: "Commands to control Testflight.",
}

// Testers

var testersCmd = &cobra.Command{
	Use:   "testers",
	Short: "Lists testers",
}

var testersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists testers",
	Run:   testersList,
}

var testersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tester",
	Run:   testersCreate,
}

var (
	providerID int
	appID      int
)

func init() {
	RootCmd.AddCommand(testflightCmd)

	testflightCmd.PersistentFlags().IntVarP(&providerID, "providerID", "p", 0, "Provider ID")
	testflightCmd.PersistentFlags().IntVarP(&appID, "appID", "a", 0, "App ID")

	testflightCmd.AddCommand(testersCmd)

	testersCmd.AddCommand(testersListCmd)
	testersCmd.AddCommand(testersCreateCmd)
}

func testersList(cmd *cobra.Command, args []string) {
	testers, err := client.TestersList(providerID, appID, nil)
	checkErr(err)
	fmt.Printf(prettyPrint(testers))
}

func testersCreate(cmd *cobra.Command, args []string) {
	fmt.Println("testerCreate: Not Implemented")
}
