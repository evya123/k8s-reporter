// utils/excel_writer.go

package utils

import (
	"os"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// OpenOrCreateExcelFile checks if an Excel file exists and opens it, or creates a new file if it doesn't exist.
func OpenOrCreateExcelFile(filePath string) (*excelize.File, error) {
	var f *excelize.File
	_, err := os.Stat(filePath)
	if err == nil {
		// File exists, open it
		f, err = excelize.OpenFile(filePath)
		if err != nil {
			Error("Failed to open existing Excel file", zap.String("filePath", filePath), zap.Error(err))
			return nil, err
		}
	} else if os.IsNotExist(err) {
		// File does not exist, create a new file
		f = excelize.NewFile()
	} else {
		// Some other error occurred when checking file existence
		Error("Failed to check Excel file existence", zap.String("filePath", filePath), zap.Error(err))
		return nil, err
	}
	return f, nil
}

// AddSheetToExcelFile adds a new sheet to the Excel file with the given headers.
func AddSheetToExcelFile(f *excelize.File, sheetName string, headers []string) error {
	// Check if the sheet already exists
	idx, err := f.GetSheetIndex(sheetName)
	if err != nil {
		Error("Could not check if sheet already exists", zap.String("sheetName", sheetName), zap.Error(err))
		return err
	} else if idx != -1 {
		Info("Sheet already exists, no need to add it again", zap.String("sheetName", sheetName))
		return nil
	}

	// Add a new sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		Error("Could not create sheet", zap.String("sheetName", sheetName), zap.Error(err))
		return err
	}
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1) // 1-based indexing
		if err != nil {
			Error("Failed to convert coordinates to cell name", zap.String("sheetName", sheetName), zap.Error(err))
			return err
		}
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			Error("Failed to set cell value", zap.String("cell", cell), zap.String("sheetName", sheetName), zap.Error(err))
			return err
		}
	}
	f.SetActiveSheet(index)
	Info("Added new sheet to Excel file", zap.String("sheetName", sheetName))

	// Delete default sheet if it exists and isn't the one we just added
	err = f.DeleteSheet("Sheet1")

	if err != nil {
		Error("Could not delete default sheet", zap.Error(err))
		return err
	}

	return nil
}
