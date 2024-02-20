/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"k8s-reporter/cmd"
	"k8s-reporter/utils"
	"log"

	"github.com/xuri/excelize/v2"
)

// Package-level variable to hold the Excel file
var excelFile *excelize.File

func main() {
	// Ensure the finalization function runs when the main function exits
	defer finalize()

	cmd.Execute()
}

// finalize is the finalization function that saves the Excel file.
func finalize() {
	log.Println("INFO: Finalizing and saving the Excel report")
	excelManager := utils.GetExcelFileManager()
	if err := excelManager.SaveExcelFile("k8s_report.xlsx"); err != nil {
		log.Fatalf("ERROR: Failed to save the Excel report: %s", err)
	}
	log.Println("INFO: Excel report saved successfully.")
}
