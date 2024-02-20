// handlers/Statefulset_handler.go

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
	"Image Versions",
	"QoS Class",
}

// FetchResources fetches all Statefulsets across all namespaces and stores them.
func (d *StatefulsetHandler) FetchResources(clientset *kubernetes.Clientset) error {
	log.Println("INFO: Fetching Statefulsets from Kubernetes cluster")
	Statefulsets, err := clientset.AppsV1().StatefulSets("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("ERROR: Failed to fetch Statefulsets: %s\n", err)
		return err
	}
	d.Statefulsets = Statefulsets.Items
	log.Printf("INFO: Fetched %d Statefulsets\n", len(d.Statefulsets))
	return nil
}

// WriteExcel writes the information of the fetched Statefulsets to an Excel sheet.
func (d *StatefulsetHandler) WriteExcel(f *excelize.File, sheetName string) error {
	log.Printf("INFO: Writing Statefulsets data to Excel sheet: %s\n", sheetName)
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, Statefulset := range d.Statefulsets {
		// Extract the necessary data from each Statefulset
		name := Statefulset.Name
		namespace := Statefulset.Namespace
		desiredReplicas := Statefulset.Spec.Replicas
		currentReplicas := Statefulset.Status.Replicas
		availableReplicas := Statefulset.Status.AvailableReplicas
		readyReplicas := Statefulset.Status.ReadyReplicas
		uptodateReplicas := Statefulset.Status.UpdatedReplicas
		nodeSelector := Statefulset.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits := utils.ExtractResources(Statefulset.Spec.Template.Spec)
		imageVersions := utils.ExtractImageVersions(Statefulset.Spec.Template.Spec)
		qosClass := utils.DetermineQoSClass(Statefulset.Spec.Template.Spec)
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
				log.Printf("ERROR: Failed to convert coordinates to cell name for Statefulset %s: %s\n", name, err)
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Printf("ERROR: Failed to set cell value for Statefulset %s: %s\n", name, err)
				return err
			}
		}
		rowIndex++
	}
	log.Printf("INFO: Successfully written Statefulset data to Excel sheet: %s\n", sheetName)
	return nil
}
