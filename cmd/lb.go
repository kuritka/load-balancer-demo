package cmd

import (
	"github.com/spf13/cobra"
	"lb/common/env"
	guards "lb/common/guard"
)

var lbCmd = &cobra.Command{
	Use:   "lb",
	Short: "load balancer",
	Long: `load balancer`,

	Run: func(cmd *cobra.Command, args []string) {

		port := env.MustGetStringFlagFromEnv("LB_PORT")
		err := lb.NewService(port).Run()

		guards.FailOnError(err)
	},
}

func init(){
	rootCmd.AddCommand(lbCmd)
}