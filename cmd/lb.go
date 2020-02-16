package cmd

import (
	"github.com/spf13/cobra"
	"lb/services/loadbalancer"

	"lb/common/env"
)

var lbCmd = &cobra.Command{
	Use:   "lb",
	Short: "load balancer",
	Long:  `load balancer`,

	Run: func(cmd *cobra.Command, args []string) {

		lbPort := env.MustGetStringFlagFromEnv("LB_PORT")
		discoPort := env.MustGetStringFlagFromEnv("DISCO_PORT")
		//i.e. ":2000", ":2001"
		loadbalancer.Run(lbPort, discoPort)
	},
}

func init() {
	rootCmd.AddCommand(lbCmd)
}
