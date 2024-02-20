# Handlers Directory

## Overview
The `handlers` directory contains structs and methods for interacting with Kubernetes resources. Each handler is responsible for fetching and writing data for a specific resource type to an Excel sheet.

## Handlers
- `daemonset_handler.go`: Handler for DaemonSets.
- `deployment_handler.go`: Handler for Deployments.
- `job_handler.go`: Handler for Jobs.
- `statefulset_handler.go`: Handler for StatefulSets.

## ResourceHandler Interface
The `handler.go` file defines the `ResourceHandler` interface, which includes the following methods:
- `FetchResources(clientset *kubernetes.Clientset) error`: Fetches resources from the Kubernetes cluster.
- `WriteCSV(writer *csv.Writer) error`: Writes resource data to a CSV file.

## Headers
Each handler file contains a `Headers` variable that defines the column headers for the Excel sheet corresponding to the resource type.

## Usage
Handlers are utilized by the commands defined in the `cmd` directory to perform resource-specific operations.
