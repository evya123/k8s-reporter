// handlers/daemonset_handler.go

package handlers

import (
	"context"
	"strconv"

	"k8s-reporter/utils"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DaemonSetHandler is a struct that implements the ResourceHandler interface
// for Kubernetes DaemonSets.
type DaemonSetHandler struct {
	DaemonSets []v1.DaemonSet
}

var DaemonSetHeaders = []string{
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

// FetchResources fetches all DaemonSets across all namespaces and stores them.
func (d *DaemonSetHandler) FetchResources(clientset *kubernetes.Clientset) error {
	utils.Info("Fetching DaemonSets from Kubernetes cluster")
	daemonSets, err := clientset.AppsV1().DaemonSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utils.Error("Failed to fetch DaemonSets", zap.Error(err))
		return err
	}
	d.DaemonSets = daemonSets.Items
	utils.Info("Fetched DaemonSets", zap.Int("count", len(d.DaemonSets)))
	return nil
}

// WriteExcel writes the information of the fetched DaemonSets to an Excel sheet.
func (d *DaemonSetHandler) WriteExcel(clientset *kubernetes.Clientset, f *excelize.File, sheetName string) error {
	utils.Info("Writing DaemonSets data to Excel sheet", zap.String("sheetName", sheetName))
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, ds := range d.DaemonSets {
		name := ds.Name
		namespace := ds.Namespace
		podSpec := ds.Spec.Template.Spec
		cpuRequests, memoryRequests, cpuLimits, memoryLimits, cpuDiff, memoryDiff, memoryReadiness, qosClass := utils.ExtractResources(clientset, podSpec, namespace)
		imageVersions := utils.ExtractImageVersions(podSpec)
		// qosClass := utils.DetermineQoSClass(podSpec)

		record := []interface{}{
			name,
			namespace,
			strconv.Itoa(int(ds.Status.DesiredNumberScheduled)),
			strconv.Itoa(int(ds.Status.CurrentNumberScheduled)),
			strconv.Itoa(int(ds.Status.NumberReady)),
			strconv.Itoa(int(ds.Status.UpdatedNumberScheduled)),
			strconv.Itoa(int(ds.Status.NumberAvailable)),
			utils.FormatNodeSelector(podSpec.NodeSelector),
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
				utils.Error("Failed to convert coordinates to cell name for DaemonSet", zap.String("daemonSetName", ds.Name), zap.Error(err))
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				utils.Error("Failed to set cell value for DaemonSet", zap.String("daemonSetName", ds.Name), zap.Error(err))
				return err
			}
		}
		rowIndex++
	}
	utils.Info("Successfully written DaemonSet data to Excel sheet", zap.String("sheetName", sheetName))
	return nil
}
