package cmd

import (
	"github.com/spf13/cobra"
	"lb/services/cache"
)

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "cache service",
	Long:  `cache service`,

	Run: func(cmd *cobra.Command, args []string) {
		cache.Run(":3080")
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
}
