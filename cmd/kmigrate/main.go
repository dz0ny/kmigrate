package main

import (
	"fmt"
	"kmigrate/logger"
	"os"

	"github.com/spf13/cobra"
)

var log = logger.New("kmigrate")

var patchPath string
var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "kmigrate",
	Short: "kmigrate is declarative patcher fo Kubernetes resources",
}

func main() {
	rootCmd.PersistentFlags().Bool("verbose", false, "Show debugging information")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
