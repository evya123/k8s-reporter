# CMD Directory

## Overview
The `cmd` directory contains the command-line interface (CLI) definitions for the `k8s-reporter` tool. Each file defines a command that allows users to export data about specific Kubernetes resources to an Excel sheet.

## Commands
- `daemonsets.go`: Export DaemonSets to an Excel sheet.
- `deployments.go`: Export Deployments to an Excel sheet.
- `jobs.go`: Export Jobs to an Excel sheet.
- `root.go`: The root command that all other commands are attached to.
- `run-all.go`: Execute all resource commands sequentially.
- `statefulsets.go`: Export StatefulSets to an Excel sheet.

## Usage
Each command can be used by running `k8s-reporter` followed by the command name.
Commands may accept a `--kubeconfig` flag to specify the path to the kubeconfig file.
For example:
`go run main.go run-all --kubeconfig=/path/to/kubeconfig` 
or
`k8s-reporter run-all --kubeconfig=/path/to/kubeconfig`