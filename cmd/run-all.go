// cmd/run-all.go

package cmd

import (
	"k8s-reporter/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:   "run-all",
	Short: "Run all resource commands",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Info("Running all resources...")
		resources := []string{"deployments", "daemonsets", "statefulsets", "jobs"}
		for _, resource := range resources {
			utils.Info("Running resource command", zap.String("resource", resource))
			rootCmd.SetArgs([]string{resource})
			if err := rootCmd.Execute(); err != nil {
				utils.Error("Failed to execute resource command", zap.String("resource", resource), zap.Error(err))
				break // Stop executing further commands after an error
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
