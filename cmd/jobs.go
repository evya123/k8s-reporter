// cmd/jobs.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// jobsCmd represents the jobs command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Export Jobs to an Excel sheet",
	Long:  `Export Jobs to an Excel sheet will fetch all the Jobs from a Kubernetes cluster and write their details to an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		utils.Info("Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			utils.Fatal("Error building Kubernetes clientset", zap.Error(err))
		}

		utils.Info("Fetching Jobs")
		jobHandler := &handlers.JobHandler{}
		if err := jobHandler.FetchResources(clientset); err != nil {
			utils.Fatal("Error fetching Jobs", zap.Error(err))
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			utils.Fatal("Failed to open or create Excel file", zap.Error(err))
		}

		excelFile := excelManager.GetExcelFile()
		utils.Info("Adding Jobs sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "Jobs", handlers.JobHeaders); err != nil {
			utils.Fatal("Failed to add sheet to Excel file", zap.Error(err))
		}

		utils.Info("Writing Jobs data to Excel sheet")
		if err := jobHandler.WriteExcel(clientset, excelFile, "Jobs"); err != nil {
			utils.Fatal("Error writing to Excel", zap.Error(err))
		}

		utils.Info("Jobs data written to Excel file successfully.")
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
