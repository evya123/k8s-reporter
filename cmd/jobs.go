// cmd/jobs.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"
	"log"

	"github.com/spf13/cobra"
)

// deploymentsCmd represents the deployments command
var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "Export Jobs to an Excel sheet",
	Long:  `Export Jobs to an Excel sheet will fetch all the Jobs from a Kubernetes cluster and write their details to an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		log.Println("INFO: Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			log.Fatalf("ERROR: Error building Kubernetes clientset: %s", err.Error())
		}

		log.Println("INFO: Fetching Jobs")
		jobHandler := &handlers.JobHandler{}
		if err := jobHandler.FetchResources(clientset); err != nil {
			log.Fatalf("ERROR: Error fetching Jobs: %s", err.Error())
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			log.Fatalf("ERROR: %s", err)
		}

		excelFile := excelManager.GetExcelFile()
		log.Println("INFO: Adding Jobs sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "Jobs", handlers.JobHeaders); err != nil {
			log.Fatalf("ERROR: Failed to add sheet to Excel file: %s", err)
		}

		log.Println("INFO: Writing Jobs data to Excel sheet")
		if err := jobHandler.WriteExcel(excelFile, "Jobs"); err != nil {
			log.Fatalf("ERROR: Error writing to Excel: %s", err)
		}

		log.Println("INFO: Jobs data written to Excel file successfully.")
	},
}

func init() {
	rootCmd.AddCommand(jobsCmd)
	jobsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
