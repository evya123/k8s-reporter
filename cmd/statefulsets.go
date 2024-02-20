// cmd/statefulsets.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// statefulsetsCmd represents the statefulsets command
var statefulsetsCmd = &cobra.Command{
	Use:   "statefulsets",
	Short: "Export Statefulsets to an Excel sheet",
	Long:  `Export Statefulsets to an Excel sheet will fetch all the Statefulsets from a Kubernetes cluster and write their details to an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		utils.Info("Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			utils.Fatal("Error building Kubernetes clientset", zap.Error(err))
		}

		utils.Info("Fetching Statefulsets")
		statefulsetHandler := &handlers.StatefulsetHandler{}
		if err := statefulsetHandler.FetchResources(clientset); err != nil {
			utils.Fatal("Error fetching Statefulsets", zap.Error(err))
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			utils.Fatal("Failed to open or create Excel file", zap.Error(err))
		}

		excelFile := excelManager.GetExcelFile()
		utils.Info("Adding Statefulsets sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "Statefulsets", handlers.StatefulsetHeaders); err != nil {
			utils.Fatal("Failed to add sheet to Excel file", zap.Error(err))
		}

		utils.Info("Writing Statefulsets data to Excel sheet")
		if err := statefulsetHandler.WriteExcel(clientset, excelFile, "Statefulsets"); err != nil {
			utils.Fatal("Error writing to Excel", zap.Error(err))
		}

		utils.Info("Statefulsets data written to Excel file successfully")
	},
}

func init() {
	rootCmd.AddCommand(statefulsetsCmd)
	statefulsetsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
