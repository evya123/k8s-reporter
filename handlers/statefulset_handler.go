// handlers/statefulset_handler.go

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

// StatefulsetHandler is a struct that implements the ResourceHandler interface
// for Kubernetes Statefulsets.
type StatefulsetHandler struct {
	Statefulsets []appsv1.StatefulSet
}

var StatefulsetHeaders = []string{
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
	"CPU Diff",
	"Memory Diff",
	"Memory diff > 2 x Request",
	"Image Versions",
	"QoS Class",
	"Owner",
}

// FetchResources fetches all Statefulsets across all namespaces and stores them.
func (d *StatefulsetHandler) FetchResources(clientset *kubernetes.Clientset) error {
	utils.Info("Fetching Statefulsets from Kubernetes cluster")
	statefulsets, err := clientset.AppsV1().StatefulSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utils.Error("Failed to fetch Statefulsets", zap.Error(err))
		return err
	}
	d.Statefulsets = statefulsets.Items
	utils.Info("Fetched Statefulsets", zap.Int("count", len(d.Statefulsets)))
	return nil
}

// WriteExcel writes the information of the fetched Statefulsets to an Excel sheet.
func (d *StatefulsetHandler) WriteExcel(clientset *kubernetes.Clientset, f *excelize.File, sheetName string) error {
	utils.Info("Writing Statefulsets data to Excel sheet", zap.String("sheetName", sheetName))
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, statefulset := range d.Statefulsets {
		name := statefulset.Name
		namespace := statefulset.Namespace
		desiredReplicas := statefulset.Spec.Replicas
		currentReplicas := statefulset.Status.Replicas
		availableReplicas := statefulset.Status.AvailableReplicas
		readyReplicas := statefulset.Status.ReadyReplicas
		uptodateReplicas := statefulset.Status.UpdatedReplicas
		nodeSelector := statefulset.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits, cpuDiff, memoryDiff, memoryReadiness, qosClass := utils.ExtractResources(clientset, statefulset.Spec.Template.Spec, namespace)
		imageVersions := utils.ExtractImageVersions(statefulset.Spec.Template.Spec)
		// qosClass := utils.DetermineQoSClass(statefulset.Spec.Template.Spec)
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
			cpuDiff,
			memoryDiff,
			memoryReadiness,
			imageVersions,
			qosClass,
		}

		for i, value := range record {
			cell, err := excelize.CoordinatesToCellName(i+1, rowIndex)
			if err != nil {
				utils.Error("Failed to convert coordinates to cell name for Statefulset", zap.String("statefulsetName", name), zap.Error(err))
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				utils.Error("Failed to set cell value for Statefulset", zap.String("statefulsetName", name), zap.Error(err))
				return err
			}
		}
		rowIndex++
	}
	utils.Info("Successfully written Statefulset data to Excel sheet", zap.String("sheetName", sheetName))
	return nil
}
