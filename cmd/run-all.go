// cmd/run-all.go

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run-all",
	Short: "Run all resource commands",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running all resources...")
		rootCmd.SetArgs([]string{"deployments"})
		rootCmd.Execute()
		rootCmd.SetArgs([]string{"daemonsets"})
		rootCmd.Execute()
		rootCmd.SetArgs([]string{"statefulsets"})
		rootCmd.Execute()
		rootCmd.SetArgs([]string{"jobs"})
		rootCmd.Execute()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
