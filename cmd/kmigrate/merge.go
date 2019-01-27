package main

import (
	"kmigrate/migrate"
	"kmigrate/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(&patchPath, "filename", "f", "", "Filename declaring patch for the resources")
	migrateCmd.Flags().BoolVarP(&dryRun, "dry-run", "", false, "If true, only print the change that would be sent, without sending it.")
}

var migrateCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge resources defined in patch file",
	Run: func(cmd *cobra.Command, args []string) {
		client := utils.GetClient()
		migrate.NewMigrate(patchPath, dryRun, client).Run()
	},
}
