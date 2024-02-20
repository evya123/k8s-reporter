/*
Copyright © 2024 Evyatar Shtern
*/
package main

import (
	"k8s-reporter/cmd"
	"k8s-reporter/utils"

	"go.uber.org/zap"
)

func main() {
	// Ensure the finalization function runs when the main function exits
	defer finalize()

	cmd.Execute()
}

// finalize is the finalization function that saves the Excel file.
func finalize() {
	utils.Info("Finalizing and saving the Excel report")
	excelManager := utils.GetExcelFileManager()
	if err := excelManager.SaveExcelFile("k8s_report.xlsx"); err != nil {
		utils.Fatal("Failed to save the Excel report: ", zap.Error(err))
	}
	utils.Info("Excel report saved successfully.")
}
