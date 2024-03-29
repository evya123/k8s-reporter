// cmd/root.go

package cmd

import (
	"k8s-reporter/utils"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "k8s-reporter",
	Short: "k8s-reporter is a CLI for creating a report about Kubernetes objects",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		utils.Error("Execution failed", zap.Error(err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("kubeconfig", "", "Path to the kubeconfig file")
}
