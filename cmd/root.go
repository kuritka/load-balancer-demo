package cmd

import (
	"os"

	"lb/common/log"

	"github.com/spf13/cobra"

)

var logger = log.Log

var Verbose bool

var rootCmd = &cobra.Command{
	Short: "load balancer demo",
	Long: `load balancer demo`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Error().Msg("No parameters included")
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logger.Info().Msg("done..")
	},
}


func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
