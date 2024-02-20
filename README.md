# k8s-reporter

`k8s-reporter` is a Command Line Interface (CLI) tool designed to generate reports about Kubernetes resources. It supports exporting information about various Kubernetes objects such as DaemonSets, Deployments, Jobs, and StatefulSets into an Excel format, providing insights into resource utilization, configuration, and status.

## Features

- Fetch and list Kubernetes resources across all namespaces in a cluster.
- Export resource details to an Excel file, including:
  - Name, Namespace, Desired Number of Pods, Current Number of Pods, Number of Ready Pods, Up-to-date Pods, Available Pods
  - Node Selector, CPU and Memory Requests, CPU and Memory Limits, Image Versions, QoS Class
- Supports multiple resource types with the ability to extend functionality for additional Kubernetes objects.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Go version 1.15 or above.
- You have a Kubernetes cluster running and have access to it.
- You have configured `kubectl` and have the appropriate context and permissions to interact with your Kubernetes cluster.
- The CLI assume that you have kubeconfig at the default location

## Installation

To install `k8s-reporter`, follow these steps:

```bash
# Clone the repository
git clone https://github.com/yourusername/k8s-reporter.git

# Navigate to the cloned directory
cd k8s-reporter

# Build the binary
go build -o k8s-reporter .
```
# Usage
Run the k8s-reporter followed by the command for the resource you want to report on. Available commands are:

* daemonsets: Export DaemonSets to an Excel sheet.
* deployments: Export Deployments to an Excel sheet.
* jobs: Export Jobs to an Excel sheet.
* statefulsets: Export StatefulSets to an Excel sheet.
* run-all: Execute all resource commands sequentially.

Example usage:
```
./k8s-reporter daemonsets
./k8s-reporter deployments --kubeconfig=/path/to/kubeconfig
./k8s-reporter run-all
```

The tool will generate a file named k8s_report.xlsx with the exported data.

Contributing
