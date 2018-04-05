package cmd

import (
	"fmt"

	"github.com/nathanborror/itc/itunes"

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
	Short: "List testers",
	Run:   testersList,
}

var testersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tester",
	Run:   testersCreate,
}

var (
	providerID   int
	appID        int
	groupID      string
	createTester itunes.CreateTester
)

func init() {
	RootCmd.AddCommand(testflightCmd)

	// testflight (required flags: --providerID, --appID)
	testflightCmd.PersistentFlags().IntVarP(&providerID, "providerID", "p", 0, "Provider ID")
	testflightCmd.PersistentFlags().IntVarP(&appID, "appID", "a", 0, "App ID")
	testflightCmd.PersistentFlags().StringVarP(&groupID, "groupID", "g", "", "Group ID")

	testflightCmd.MarkPersistentFlagRequired("providerID")
	testflightCmd.MarkPersistentFlagRequired("appID")

	// testflight testers (noop)
	testflightCmd.AddCommand(testersCmd)

	// testflight testers list
	testersCmd.AddCommand(testersListCmd)

	// testflight testers create (required flags: --email, --first-name, --last-name, --groupID)
	testersCreateCmd.Flags().StringVar(&createTester.Email, "email", "", "Tester email")
	testersCreateCmd.Flags().StringVar(&createTester.FirstName, "first-name", "", "Tester first name")
	testersCreateCmd.Flags().StringVar(&createTester.LastName, "last-name", "", "Tester last name")
	testersCreateCmd.MarkFlagRequired("email")
	testersCreateCmd.MarkFlagRequired("first-name")
	testersCreateCmd.MarkFlagRequired("last-name")
	testersCmd.AddCommand(testersCreateCmd)
}

func testersList(cmd *cobra.Command, args []string) {
	testers, err := client.TestersList(providerID, appID, nil)
	checkErr(err)
	fmt.Printf(prettyPrint(testers))
}

func testersCreate(cmd *cobra.Command, args []string) {
	err := client.TesterCreate([]itunes.CreateTester{createTester}, providerID, appID, groupID)
	checkErr(err)
}
