// cmd/deployments.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"
	"log"

	"github.com/spf13/cobra"
)

// deploymentsCmd represents the deployments command
var deploymentsCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Export Deployments to an Excel sheet",
	Long:  `Export Deployments to an Excel sheet will fetch all the Deployments from a Kubernetes cluster and write their details to an Excel file.`,
	RunE:  deployments,
}

func deployments(cmd *cobra.Command, args []string) error {
	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")
	log.Println("INFO: Building Kubernetes clientset")
	clientset, err := utils.GetKubernetesClient(kubeconfig)
	if err != nil {
		log.Fatalf("ERROR: Error building Kubernetes clientset: %s", err.Error())
		return err
	}

	log.Println("INFO: Fetching Deployments")
	deploymentHandler := &handlers.DeploymentHandler{}
	if err := deploymentHandler.FetchResources(clientset); err != nil {
		log.Fatalf("ERROR: Error fetching Deployments: %s", err.Error())
		return err
	}

	excelManager := utils.GetExcelFileManager()
	if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
		log.Fatalf("ERROR: %s", err)
		return err
	}

	excelFile := excelManager.GetExcelFile()
	log.Println("INFO: Adding Deployments sheet to Excel file")
	if err := utils.AddSheetToExcelFile(excelFile, "Deployments", handlers.DeploymentHeaders); err != nil {
		log.Fatalf("ERROR: Failed to add sheet to Excel file: %s", err)
		return err
	}

	log.Println("INFO: Writing Deployments data to Excel sheet")
	if err := deploymentHandler.WriteExcel(excelFile, "Deployments"); err != nil {
		log.Fatalf("ERROR: Error writing to Excel: %s", err)
		return err
	}

	log.Println("INFO: Deployments data written to Excel file successfully.")
	return nil
}

func init() {
	rootCmd.AddCommand(deploymentsCmd)
	deploymentsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
