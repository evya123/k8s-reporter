// cmd/daemonsets.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// daemonsetsCmd represents the daemonsets command
var daemonsetsCmd = &cobra.Command{
	Use:   "daemonsets",
	Short: "Export DaemonSets to an Excel sheet",
	Long: `Export DaemonSets to an Excel sheet will fetch all the DaemonSets from a Kubernetes cluster
and write their details to a specified sheet within an Excel file.`,
	Example: `# Export DaemonSets to an Excel sheet using the default kubeconfig
k8s-reporter daemonsets

# Export DaemonSets to an Excel sheet using a specific kubeconfig
k8s-reporter daemonsets --kubeconfig=/path/to/kubeconfig`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		utils.Info("Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			utils.Fatal("Error building Kubernetes clientset", zap.Error(err))
		}

		daemonSetHandler := &handlers.DaemonSetHandler{}
		utils.Info("Fetching DaemonSets")
		if err := daemonSetHandler.FetchResources(clientset); err != nil {
			utils.Fatal("Error fetching DaemonSets", zap.Error(err))
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			utils.Fatal("Failed to open or create Excel file", zap.Error(err))
		}

		excelFile := excelManager.GetExcelFile()

		utils.Info("Adding DaemonSets sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "DaemonSets", handlers.DaemonSetHeaders); err != nil {
			utils.Fatal("Failed to add sheet to Excel file", zap.Error(err))
		}

		utils.Info("Writing DaemonSets data to Excel sheet")
		if err := daemonSetHandler.WriteExcel(clientset, excelFile, "DaemonSets"); err != nil {
			utils.Fatal("Error writing to Excel", zap.Error(err))
		}

		utils.Info("DaemonSets data written to Excel file successfully")
	},
}

func init() {
	rootCmd.AddCommand(daemonsetsCmd)
	daemonsetsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
