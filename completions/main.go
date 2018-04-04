package main

import (
	"fmt"
	"os"

	"github.com/nathanborror/itc/cmd"
)

func main() {
	err := cmd.RootCmd.GenBashCompletionFile("completions.sh")
	if err != nil {
		fmt.Printf("Failed to generate bash completions: %v\n", err)
		os.Exit(1)
	}
}
