// utils/namespace_info.go

package utils

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// getLimitRangeItems fetches the LimitRange items for a given namespace using an existing clientset.
func getLimitRangeItems(clientset *kubernetes.Clientset, namespace string) ([]v1.LimitRangeItem, error) {
	limitRanges, err := clientset.CoreV1().LimitRanges(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			Info("No LimitRange found in namespace", zap.String("namespace", namespace))
		} else {
			Error("Error getting LimitRange for namespace", zap.String("namespace", namespace), zap.Error(err))
		}
		return nil, err
	}

	if len(limitRanges.Items) == 0 {
		return nil, fmt.Errorf("no LimitRange items found in namespace: %s", namespace)
	}

	return limitRanges.Items[0].Spec.Limits, nil
}

// GetNamespaceDefaultCPURequests retrieves the default CPU request for a namespace using an existing clientset.
func GetNamespaceDefaultCPURequests(clientset *kubernetes.Clientset, namespace string) resource.Quantity {
	limitItems, err := getLimitRangeItems(clientset, namespace)
	if err != nil {
		return resource.Quantity{}
	}

	for _, item := range limitItems {
		if item.Type == v1.LimitTypeContainer {
			if val, ok := item.DefaultRequest[v1.ResourceCPU]; ok {
				return val
			}
		}
	}

	Info("Default CPU request not found in namespace", zap.String("namespace", namespace))
	return resource.Quantity{}
}

// GetNamespaceDefaultMemoryRequests retrieves the default Memory request for a namespace using an existing clientset.
func GetNamespaceDefaultMemoryRequests(clientset *kubernetes.Clientset, namespace string) resource.Quantity {
	limitItems, err := getLimitRangeItems(clientset, namespace)
	if err != nil {
		return resource.Quantity{}
	}

	for _, item := range limitItems {
		if item.Type == v1.LimitTypeContainer {
			if val, ok := item.DefaultRequest[v1.ResourceMemory]; ok {
				return val
			}
		}
	}

	Info("Default Memory request not found in namespace", zap.String("namespace", namespace))
	return resource.Quantity{}
}

// GetNamespaceDefaultCPULimits retrieves the default CPU limit for a namespace using an existing clientset.
func GetNamespaceDefaultCPULimits(clientset *kubernetes.Clientset, namespace string) resource.Quantity {
	limitItems, err := getLimitRangeItems(clientset, namespace)
	if err != nil {
		return resource.Quantity{}
	}

	for _, item := range limitItems {
		if item.Type == v1.LimitTypeContainer {
			if val, ok := item.Default[v1.ResourceCPU]; ok {
				return val
			}
		}
	}

	Info("Default CPU limit not found in namespace", zap.String("namespace", namespace))
	return resource.Quantity{}
}

// GetNamespaceDefaultMemoryLimits retrieves the default Memory limit for a namespace using an existing clientset.
func GetNamespaceDefaultMemoryLimits(clientset *kubernetes.Clientset, namespace string) resource.Quantity {
	limitItems, err := getLimitRangeItems(clientset, namespace)
	if err != nil {
		return resource.Quantity{}
	}

	for _, item := range limitItems {
		if item.Type == v1.LimitTypeContainer {
			if val, ok := item.Default[v1.ResourceMemory]; ok {
				return val
			}
		}
	}

	Info("Default Memory limit not found in namespace", zap.String("namespace", namespace))
	return resource.Quantity{}
}
