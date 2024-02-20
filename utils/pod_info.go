// utils/pod_info.go

package utils

import (
	"fmt"
	"strconv"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// formatNodeSelector takes a map and returns a formatted string of key-value pairs
func FormatNodeSelector(nodeSelector map[string]string) string {
	var selectorPairs []string
	for key, value := range nodeSelector {
		selectorPairs = append(selectorPairs, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(selectorPairs, ", ")
}

// FormatResourceQuantity converts a resource.Quantity to a string and ensures memory is in MiB.
func FormatResourceQuantity(q resource.Quantity, resourceName v1.ResourceName) string {
	if resourceName == v1.ResourceMemory {
		// If the quantity is already in MiB, return it as is
		if strings.HasSuffix(q.String(), "Mi") {
			return q.String()
		}
		// Otherwise, convert to MiB
		return strconv.FormatInt(q.ScaledValue(resource.Mega), 10) + "Mi"
	}
	return q.String()
}

// ExtractResources takes a PodSpec and returns formatted strings of CPU and memory requests and limits.
func ExtractResources(podSpec v1.PodSpec) (cpuRequests, memoryRequests, cpuLimits, memoryLimits string) {
	var cpuReqTotal, memReqTotal, cpuLimitTotal, memLimitTotal resource.Quantity
	for _, container := range podSpec.Containers {
		if cpu, ok := container.Resources.Requests[v1.ResourceCPU]; ok {
			cpuReqTotal.Add(cpu)
		}
		if memory, ok := container.Resources.Requests[v1.ResourceMemory]; ok {
			memReqTotal.Add(memory)
		}
		if cpu, ok := container.Resources.Limits[v1.ResourceCPU]; ok {
			cpuLimitTotal.Add(cpu)
		}
		if memory, ok := container.Resources.Limits[v1.ResourceMemory]; ok {
			memLimitTotal.Add(memory)
		}
	}

	cpuRequests = strconv.FormatInt(cpuReqTotal.MilliValue(), 10) + "m"
	memoryRequests = FormatResourceQuantity(memReqTotal, v1.ResourceMemory)
	cpuLimits = strconv.FormatInt(cpuLimitTotal.MilliValue(), 10) + "m"
	memoryLimits = FormatResourceQuantity(memLimitTotal, v1.ResourceMemory)

	return
}

// ExtractImageVersions takes a PodSpec and returns a string containing image versions (tags).
func ExtractImageVersions(podSpec v1.PodSpec) string {
	var imageVersions []string
	for _, container := range podSpec.Containers {
		image := container.Image
		parts := strings.Split(image, ":")
		var version string
		if len(parts) > 1 {
			version = parts[1]
		} else {
			version = "latest" // Default tag if none is specified
		}
		imageVersions = append(imageVersions, version)
	}
	return strings.Join(imageVersions, ", ")
}

// DetermineQoSClass takes a PodSpec and returns its QoS class as a string.
func DetermineQoSClass(podSpec v1.PodSpec) string {
	guaranteed := true
	burstable := false
	for _, container := range podSpec.Containers {
		requests := container.Resources.Requests
		limits := container.Resources.Limits
		if len(requests) == 0 && len(limits) == 0 {
			guaranteed = false
			continue
		}

		if _, ok := requests[v1.ResourceCPU]; ok {
			burstable = true
		}
		if _, ok := requests[v1.ResourceMemory]; ok {
			burstable = true
		}

		if _, ok := limits[v1.ResourceCPU]; !ok {
			guaranteed = false
		}
		if _, ok := limits[v1.ResourceMemory]; !ok {
			guaranteed = false
		}

		if cpuLimit, ok := limits[v1.ResourceCPU]; ok {
			if cpuRequest, ok := requests[v1.ResourceCPU]; ok {
				if cpuLimit.Cmp(cpuRequest) != 0 {
					guaranteed = false
				}
			} else {
				guaranteed = false
			}
		}

		if memLimit, ok := limits[v1.ResourceMemory]; ok {
			if memRequest, ok := requests[v1.ResourceMemory]; ok {
				if memLimit.Cmp(memRequest) != 0 {
					guaranteed = false
				}
			} else {
				guaranteed = false
			}
		}
	}

	if guaranteed {
		return "Guaranteed"
	}
	if burstable {
		return "Burstable"
	}
	return "BestEffort"
}
