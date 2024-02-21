// handlers/job_handler.go

package handlers

import (
	"context"
	"k8s-reporter/utils"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
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
	"CPU Diff",
	"Memory Diff",
	"Memory diff > 2 x Request",
	"Image Versions",
	"QoS Class",
	"Owner",
}

// FetchResources fetches all Jobs across all namespaces and stores them.
func (j *JobHandler) FetchResources(clientset *kubernetes.Clientset) error {
	utils.Info("Fetching Jobs from Kubernetes cluster")
	jobs, err := clientset.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		utils.Error("Failed to fetch Jobs", zap.Error(err))
		return err
	}
	j.Jobs = jobs.Items
	utils.Info("Fetched Jobs", zap.Int("count", len(j.Jobs)))
	return nil
}

// WriteExcel writes the information of the fetched Jobs to an Excel sheet.
func (j *JobHandler) WriteExcel(clientset *kubernetes.Clientset, f *excelize.File, sheetName string) error {
	utils.Info("Writing Jobs data to Excel sheet", zap.String("sheetName", sheetName))
	// Starting from the second row, since the first row is for headers
	rowIndex := 2
	for _, job := range j.Jobs {
		name := job.Name
		namespace := job.Namespace
		nodeSelector := job.Spec.Template.Spec.NodeSelector
		cpuRequests, memoryRequests, cpuLimits, memoryLimits, cpuDiff, memoryDiff, memoryReadiness, qosClass := utils.ExtractResources(clientset, job.Spec.Template.Spec, namespace)
		imageVersions := utils.ExtractImageVersions(job.Spec.Template.Spec)
		// qosClass := utils.DetermineQoSClass(job.Spec.Template.Spec)

		record := []interface{}{
			name,
			namespace,
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
				utils.Error("Failed to convert coordinates to cell name for Job", zap.String("jobName", name), zap.Error(err))
				return err
			}
			if err := f.SetCellValue(sheetName, cell, value); err != nil {
				utils.Error("Failed to set cell value for Job", zap.String("jobName", name), zap.Error(err))
				return err
			}
		}
		rowIndex++
	}
	utils.Info("Successfully written Job data to Excel sheet", zap.String("sheetName", sheetName))
	return nil
}
