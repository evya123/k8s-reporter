# k8s-reporter

`k8s-reporter` is a Command Line Interface (CLI) tool designed to generate reports about Kubernetes resources. It currently supports exporting information about DaemonSets within a Kubernetes cluster into a CSV format, providing insights into resource utilization, configuration, and status.

## Features

- Fetch and list all DaemonSets across all namespaces in a Kubernetes cluster.
- Export DaemonSet details to a CSV file, including:
  - Name, Namespace, Desired Number of Pods, Current Number of Pods, Number of Ready Pods, Up-to-date Pods, Available Pods
  - Node Selector, CPU and Memory Requests, CPU and Memory Limits, Image Versions, QoS Class

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Go version 1.15 or above.
- You have a Kubernetes cluster running and have access to it.
- You have configured `kubectl` and have the appropriate context and permissions to interact with your Kubernetes cluster.

## Installation

To install `k8s-reporter`, follow these steps:

```bash
# Clone the repository
git clone https://github.com/yourusername/k8s-reporter.git

# Navigate to the cloned directory
cd k8s-reporter

# Build the binary
go build -o k8s-reporter main.go
```

## Usage

To export DaemonSets to CSV using the default kubeconfig
`./k8s-reporter daemonsets`

To export DaemonSets to CSV using a specific kubeconfig
`./k8s-reporter daemonsets --kubeconfig=/path/to/kubeconfig`
The CSV file daemonsets_info.csv will be created in the current working directory with the exported data.


License
Distributed under the MIT License. See LICENSE for more information.