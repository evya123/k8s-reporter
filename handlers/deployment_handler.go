// handlers/deployment_handler.go

package handlers

import (
	"context"
	"k8s-reporter/utils"
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
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
}

// FetchResources fetches all Deployments across all namespaces and stores them.
func (d *DeploymentHandler) FetchResources(clientset *kubernetes.Clientset) error {
	log.Println("INFO: Fetching Deployments from Kubernetes cluster")
	deployments, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("ERROR: Failed to fetch Deployments: %s\n", err)
		return err
	}
	d.Deployments = deployments.Items
	log.Printf("INFO: Fetched %d Deployments\n", len(d.Deployments))
	return nil
}

// WriteExcel writes the information of the fetched Deployments to an Excel sheet.
func (d *DeploymentHandler) WriteExcel(f *excelize.File, sheetName string) error {
	log.Printf("INFO: Writing Deployments data to Excel sheet: %s\n", sheetName)
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, deployment := range d.Deployments {
		// Extract the necessary data from each Deployment
		name := deployment.Name
		namespace := deployment.Namespace
		desiredReplicas := deployment.Spec.Replicas
		currentReplicas := deployment.Status.Replicas
		availableReplicas := deployment.Status.AvailableReplicas
		readyReplicas := deployment.Status.ReadyReplicas
		uptodateReplicas := deployment.Status.UpdatedReplicas
		nodeSelector := deployment.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits := utils.ExtractResources(deployment.Spec.Template.Spec)
		imageVersions := utils.ExtractImageVersions(deployment.Spec.Template.Spec)
		qosClass := utils.DetermineQoSClass(deployment.Spec.Template.Spec)
		// Convert int32 pointers to int for display purposes
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
				log.Printf("ERROR: Failed to convert coordinates to cell name for Deployment %s: %s\n", name, err)
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Printf("ERROR: Failed to set cell value for Deployment %s: %s\n", name, err)
				return err
			}
		}
		rowIndex++
	}
	log.Printf("INFO: Successfully written Deployment data to Excel sheet: %s\n", sheetName)
	return nil
}
