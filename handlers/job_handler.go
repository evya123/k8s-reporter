// handlers/job_handler.go

package handlers

import (
	"context"
	"k8s-reporter/utils"
	"log"

	"github.com/xuri/excelize/v2"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// JobHandler is a struct that implements the ResourceHandler interface
// for Kubernetes Jobs.
type JobHandler struct {
	Jobs []batchv1.Job
}

var JobHeaders = []string{
	"Name",
	"Namespace",
	"Node Selector",
	"CPU Requests",
	"Memory Requests",
	"CPU Limits",
	"Memory Limits",
	"Image Versions",
	"QoS Class",
}

// FetchResources fetches all Jobs across all namespaces and stores them.
func (j *JobHandler) FetchResources(clientset *kubernetes.Clientset) error {
	log.Println("INFO: Fetching Jobs from Kubernetes cluster")
	jobs, err := clientset.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Printf("ERROR: Failed to fetch Jobs: %s\n", err)
		return err
	}
	j.Jobs = jobs.Items
	log.Printf("INFO: Fetched %d Jobs\n", len(j.Jobs))
	return nil
}

// WriteExcel writes the information of the fetched Jobs to an Excel sheet.
func (j *JobHandler) WriteExcel(f *excelize.File, sheetName string) error {
	log.Printf("INFO: Writing Jobs data to Excel sheet: %s\n", sheetName)
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, job := range j.Jobs {
		// Extract the necessary data from each Job
		name := job.Name
		namespace := job.Namespace
		nodeSelector := job.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits := utils.ExtractResources(job.Spec.Template.Spec)
		imageVersions := utils.ExtractImageVersions(job.Spec.Template.Spec)
		qosClass := utils.DetermineQoSClass(job.Spec.Template.Spec)

		// Convert int32 pointers to int for display purposes

		record := []interface{}{
			name,
			namespace,
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
				log.Printf("ERROR: Failed to convert coordinates to cell name for Job %s: %s\n", name, err)
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				log.Printf("ERROR: Failed to set cell value for Job %s: %s\n", name, err)
				return err
			}
		}
		rowIndex++
	}
	log.Printf("INFO: Successfully written Job data to Excel sheet: %s\n", sheetName)
	return nil
}
