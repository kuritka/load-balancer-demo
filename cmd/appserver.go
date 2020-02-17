package cmd

import (
	"lb/common/env"
	appserver "lb/server"

	"github.com/spf13/cobra"
)

var appServerCmd = &cobra.Command{
	Use:   "appserver",
	Short: "appserver",
	Long:  `appserver`,

	Run: func(cmd *cobra.Command, args []string) {

		lbDiscoUrl := env.MustGetStringFlagFromEnv("LB_DISCO_URL")
		cashServiceUrl := env.MustGetStringFlagFromEnv("CASH_SVC_URL")
		// "https://127.0.0.1:2001"
		appserver.Run(lbDiscoUrl, cashServiceUrl)
	},
}

func init() {
	rootCmd.AddCommand(appServerCmd)
}
