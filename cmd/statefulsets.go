// cmd/deployments.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"
	"log"

	"github.com/spf13/cobra"
)

// deploymentsCmd represents the deployments command
var statefulsetsCmd = &cobra.Command{
	Use:   "statefulsets",
	Short: "Export Statefulsets to an Excel sheet",
	Long:  `Export Statefulsets to an Excel sheet will fetch all the Statefulsets from a Kubernetes cluster and write their details to an Excel file.`,
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
		log.Println("INFO: Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			log.Fatalf("ERROR: Error building Kubernetes clientset: %s", err.Error())
		}

		log.Println("INFO: Fetching Statefulsets")
		statefulsetHandler := &handlers.StatefulsetHandler{}
		if err := statefulsetHandler.FetchResources(clientset); err != nil {
			log.Fatalf("ERROR: Error fetching Statefulsets: %s", err.Error())
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			log.Fatalf("ERROR: %s", err)
		}

		excelFile := excelManager.GetExcelFile()
		log.Println("INFO: Adding Statefulsets sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "Statefulsets", handlers.StatefulsetHeaders); err != nil {
			log.Fatalf("ERROR: Failed to add sheet to Excel file: %s", err)
		}

		log.Println("INFO: Writing Statefulsets data to Excel sheet")
		if err := statefulsetHandler.WriteExcel(excelFile, "Statefulsets"); err != nil {
			log.Fatalf("ERROR: Error writing to Excel: %s", err)
		}

		log.Println("INFO: Statefulsets data written to Excel file successfully.")
	},
}

func init() {
	rootCmd.AddCommand(statefulsetsCmd)
	statefulsetsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
