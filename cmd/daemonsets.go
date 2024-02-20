// cmd/daemonsets.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"
	"log"

	"github.com/spf13/cobra"
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
		log.Println("INFO: Building Kubernetes clientset")
		clientset, err := utils.GetKubernetesClient(kubeconfig)
		if err != nil {
			log.Fatalf("ERROR: Error building Kubernetes clientset: %s", err.Error())
		}

		daemonSetHandler := &handlers.DaemonSetHandler{}
		log.Println("INFO: Fetching DaemonSets")
		if err := daemonSetHandler.FetchResources(clientset); err != nil {
			log.Fatalf("ERROR: Error fetching DaemonSets: %s", err.Error())
		}

		excelManager := utils.GetExcelFileManager()
		if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
			log.Fatalf("ERROR: %s", err)
		}

		excelFile := excelManager.GetExcelFile()

		log.Println("INFO: Adding DaemonSets sheet to Excel file")
		if err := utils.AddSheetToExcelFile(excelFile, "DaemonSets", handlers.DaemonSetHeaders); err != nil {
			log.Fatalf("ERROR: Failed to add sheet to Excel file: %s", err)
		}

		log.Println("INFO: Writing DaemonSets data to Excel sheet")
		if err := daemonSetHandler.WriteExcel(excelFile, "DaemonSets"); err != nil {
			log.Fatalf("ERROR: Error writing to Excel: %s", err)
		}

		log.Println("INFO: DaemonSets data written to Excel file successfully.")
	},
}

func init() {
	rootCmd.AddCommand(daemonsetsCmd)
	daemonsetsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
