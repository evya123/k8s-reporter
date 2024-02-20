// cmd/deployments.go

package cmd

import (
	"k8s-reporter/handlers"
	"k8s-reporter/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	utils.Info("Building Kubernetes clientset")
	clientset, err := utils.GetKubernetesClient(kubeconfig)
	if err != nil {
		utils.Fatal("Error building Kubernetes clientset", zap.Error(err))
		return err
	}

	utils.Info("Fetching Deployments")
	deploymentHandler := &handlers.DeploymentHandler{}
	if err := deploymentHandler.FetchResources(clientset); err != nil {
		utils.Fatal("Error fetching Deployments", zap.Error(err))
		return err
	}

	excelManager := utils.GetExcelFileManager()
	if err := excelManager.OpenOrCreateExcelFile("k8s_report.xlsx"); err != nil {
		utils.Fatal("Failed to open or create Excel file", zap.Error(err))
		return err
	}

	excelFile := excelManager.GetExcelFile()
	utils.Info("Adding Deployments sheet to Excel file")
	if err := utils.AddSheetToExcelFile(excelFile, "Deployments", handlers.DeploymentHeaders); err != nil {
		utils.Fatal("Failed to add sheet to Excel file", zap.Error(err))
		return err
	}

	utils.Info("Writing Deployments data to Excel sheet")
	if err := deploymentHandler.WriteExcel(clientset, excelFile, "Deployments"); err != nil {
		utils.Fatal("Error writing to Excel", zap.Error(err))
		return err
	}

	utils.Info("Deployments data written to Excel file successfully")
	return nil
}

func init() {
	rootCmd.AddCommand(deploymentsCmd)
	deploymentsCmd.Flags().String("kubeconfig", "", "Path to the kubeconfig file (optional if environment variable KUBECONFIG is set)")
}
