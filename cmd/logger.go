package cmd

import (
	"github.com/spf13/cobra"
	"lb/services/logaggregator"
)

var logCmd = &cobra.Command{
	Use:   "logger",
	Short: "log aggregator",
	Long:  `log aggregator`,

	Run: func(cmd *cobra.Command, args []string) {
		logaggregator.Run(":6080")
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
