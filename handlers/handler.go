// handlers/handler.go

package handlers

import (
	"encoding/csv"

	"k8s.io/client-go/kubernetes"
)

// ResourceHandler defines the methods required to fetch Kubernetes resources
// and write their information to a CSV file.
type ResourceHandler interface {
	FetchResources(clientset *kubernetes.Clientset) error
	WriteCSV(writer *csv.Writer) error
}
