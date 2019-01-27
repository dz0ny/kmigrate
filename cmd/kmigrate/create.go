package main

import (
	"kmigrate/migrate"
	"kmigrate/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&patchPath, "filename", "f", "", "Filename declaring patch for the resources")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create patch from provided mergefile",
	Run: func(cmd *cobra.Command, args []string) {
		client := utils.GetClient()
		migrate.NewMigrate(patchPath, dryRun, client).Create()
	},
}
