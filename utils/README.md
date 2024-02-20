# Utils Directory

## Overview
The `utils` directory contains utility functions and types that provide support for Excel file manipulation, Kubernetes client initialization, and pod resource information formatting.

## Contents
- `excel_manager.go`: Manages a singleton instance of an Excel file for operations like opening, creating, and saving.
- `excel_writer.go`: Provides functions to open or create Excel files and to add new sheets with specified headers.
- `k8s_client.go`: Initializes a Kubernetes clientset using the default kubeconfig path or a specified path.
- `pod_info.go`: Includes several functions to format node selectors, resource quantities, extract resource requests and limits from pod specs, determine image versions, and identify the QoS class of a pod.

## Usage
These utilities are used throughout the `k8s-reporter` tool to facilitate interactions with Kubernetes objects and to generate reports in Excel format.