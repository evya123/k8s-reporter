# Utils Directory

## Overview
The `utils` directory contains utility functions and types that provide support for Excel file manipulation, Kubernetes client initialization, pod resource information formatting, and retrieval of default namespace resources.

## Contents
- `excel_manager.go`: Manages a singleton instance of an Excel file for operations like opening, creating, and saving.
- `excel_writer.go`: Provides functions to open or create Excel files and to add new sheets with specified headers.
- `k8s_client.go`: Initializes a Kubernetes clientset using the default kubeconfig path or a specified path.
- `pod_info.go`: Includes several functions to:
  - Format node selectors (`FormatNodeSelector`).
  - Convert and format resource quantities (`FormatResourceQuantity`).
  - Extract resource requests and limits from pod specs (`ExtractResources`).
  - Determine image versions used in a pod (`ExtractImageVersions`).
  - Identify the QoS class of a pod (`DetermineQoSClass`).
  - Retrieve default CPU and memory requests and limits for a namespace (`GetNamespaceDefaultResources`).

## Usage
These utilities are used throughout the `k8s-reporter` tool to facilitate interactions with Kubernetes objects and to generate reports in Excel format. The `GetNamespaceDefaultResources` function can be used to fetch the default resource constraints applied to pods within a specific namespace, which is particularly useful for reporting and auditing purposes.