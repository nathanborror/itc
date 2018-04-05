package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nathanborror/itc/itunes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFlag  string
	verboseFlag bool
	client      *itunes.Client
)

var RootCmd = &cobra.Command{
	Use:   "itc",
	Short: "Control your iTunes Connect account.",
}

func Execute() {
	checkErr(RootCmd.Execute())

	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "config file (default is ~/.itc.yaml)")

	RootCmd.PersistentFlags().String("appleID", "", "Your Apple ID")
	RootCmd.PersistentFlags().String("appleIDPassword", "", "Your Apple ID password")

	viper.BindPFlag("appleID", RootCmd.PersistentFlags().Lookup("appleID"))
	viper.BindPFlag("appleIDPassword", RootCmd.PersistentFlags().Lookup("appleIDPassword"))
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if configFlag != "" {
		viper.SetConfigFile(configFlag)
	}
	viper.SetConfigName(".itc")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("itc")
	err := viper.ReadInConfig()
	if err == nil && verboseFlag {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	id := viper.GetString("appleID")
	password := viper.GetString("appleIDPassword")
	client, err = itunes.NewClient(id, password)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		printErr(err.Error())
	}
}

func printJSON(in interface{}) {
	b, err := json.MarshalIndent(in, "", "  ")
	checkErr(err)
	print("%s", b)
}

func printErr(format string, a ...interface{}) {
	print(format+"\n", a...)
	os.Exit(1)
}

func print(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
