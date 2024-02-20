// handlers/deployment_handler.go

package handlers

import (
	"context"
	"k8s-reporter/utils"
	"strconv"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeploymentHandler is a struct that implements the ResourceHandler interface
// for Kubernetes Deployments.
type DeploymentHandler struct {
	Deployments []appsv1.Deployment
}

var DeploymentHeaders = []string{
	"Name",
	"Namespace",
	"Desired",
	"Current",
	"Ready",
	"Up-to-date",
	"Available",
	"Node Selector",
	"CPU Requests",
	"Memory Requests",
	"CPU Limits",
	"Memory Limits",
	"Image Versions",
	"QoS Class",
	"Owner",
}

// FetchResources fetches all Deployments across all namespaces and stores them.
func (d *DeploymentHandler) FetchResources(clientset *kubernetes.Clientset) error {
	utils.Info("Fetching Deployments from Kubernetes cluster")
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utils.Error("Failed to fetch Deployments", zap.Error(err))
		return err
	}
	d.Deployments = deployments.Items
	utils.Info("Fetched Deployments", zap.Int("count", len(d.Deployments)))
	return nil
}

// WriteExcel writes the information of the fetched Deployments to an Excel sheet.
func (d *DeploymentHandler) WriteExcel(clientset *kubernetes.Clientset, f *excelize.File, sheetName string) error {
	utils.Info("Writing Deployments data to Excel sheet", zap.String("sheetName", sheetName))
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, deployment := range d.Deployments {
		name := deployment.Name
		namespace := deployment.Namespace
		desiredReplicas := deployment.Spec.Replicas
		currentReplicas := deployment.Status.Replicas
		availableReplicas := deployment.Status.AvailableReplicas
		readyReplicas := deployment.Status.ReadyReplicas
		uptodateReplicas := deployment.Status.UpdatedReplicas
		nodeSelector := deployment.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits := utils.ExtractResources(clientset, deployment.Spec.Template.Spec, namespace)
		imageVersions := utils.ExtractImageVersions(deployment.Spec.Template.Spec)
		qosClass := utils.DetermineQoSClass(deployment.Spec.Template.Spec)
		desired := "unknown"
		if desiredReplicas != nil {
			desired = strconv.Itoa(int(*desiredReplicas))
		}

		record := []interface{}{
			name,
			namespace,
			desired,
			strconv.Itoa(int(currentReplicas)),
			strconv.Itoa(int(readyReplicas)),
			strconv.Itoa(int(uptodateReplicas)),
			strconv.Itoa(int(availableReplicas)),
			utils.FormatNodeSelector(nodeSelector),
			cpuRequests,
			memoryRequests,
			cpuLimits,
			memoryLimits,
			imageVersions,
			qosClass,
		}

		for i, value := range record {
			cell, err := excelize.CoordinatesToCellName(i+1, rowIndex)
			if err != nil {
				utils.Error("Failed to convert coordinates to cell name for Deployment", zap.String("deploymentName", name), zap.Error(err))
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				utils.Error("Failed to set cell value for Deployment", zap.String("deploymentName", name), zap.Error(err))
				return err
			}
		}
		rowIndex++
	}
	utils.Info("Successfully written Deployment data to Excel sheet", zap.String("sheetName", sheetName))
	return nil
}
