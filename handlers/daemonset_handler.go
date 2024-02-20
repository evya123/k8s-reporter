// handlers/daemonset_handler.go

package handlers

import (
	"context"
	"log"
	"strconv"

	"k8s-reporter/utils"

	"github.com/xuri/excelize/v2"
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
	"Image Versions",
	"QoS Class",
}

// FetchResources fetches all DaemonSets across all namespaces and stores them.
func (d *DaemonSetHandler) FetchResources(clientset *kubernetes.Clientset) error {
	log.Println("INFO: Fetching DaemonSets from Kubernetes cluster")
	daemonSets, err := clientset.AppsV1().DaemonSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("ERROR: Failed to fetch DaemonSets: %s\n", err)
		return err
	}
	d.DaemonSets = daemonSets.Items
	log.Printf("INFO: Fetched %d DaemonSets\n", len(d.DaemonSets))
	return nil
}

// WriteExcel writes the information of the fetched DaemonSets to an Excel sheet.
func (d *DaemonSetHandler) WriteExcel(f *excelize.File, sheetName string) error {
	log.Printf("INFO: Writing DaemonSets data to Excel sheet: %s\n", sheetName)
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, ds := range d.DaemonSets {
		podSpec := ds.Spec.Template.Spec
		cpuRequests, memoryRequests, cpuLimits, memoryLimits := utils.ExtractResources(podSpec)
		imageVersions := utils.ExtractImageVersions(podSpec)
		qosClass := utils.DetermineQoSClass(podSpec)

		record := []interface{}{
			ds.Name,
			ds.Namespace,
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
			imageVersions,
			qosClass,
		}

		for i, value := range record {
			cell, err := excelize.CoordinatesToCellName(i+1, rowIndex)
			if err != nil {
				log.Printf("ERROR: Failed to convert coordinates to cell name for DaemonSet %s: %s\n", ds.Name, err)
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Printf("ERROR: Failed to set cell value for DaemonSet %s: %s\n", ds.Name, err)
				return err
			}
		}
		rowIndex++
	}
	log.Printf("INFO: Successfully written DaemonSet data to Excel sheet: %s\n", sheetName)
	return nil
}
