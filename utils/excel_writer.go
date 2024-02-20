// utils/excel_writer.go

package utils

import (
	"os"

	"log"

	"github.com/xuri/excelize/v2"
)

// OpenOrCreateExcelFile checks if an Excel file exists and opens it, or creates a new file if it doesn't exist.
func OpenOrCreateExcelFile(filePath string) (*excelize.File, error) {
	var f *excelize.File
	_, err := os.Stat(filePath)
	if err == nil {
		// File exists, open it
		f, err = excelize.OpenFile(filePath)
		if err != nil {
			log.Printf("ERROR: Failed to open existing Excel file: %s\n", err)
			return nil, err
		}
	} else if os.IsNotExist(err) {
		// File does not exist, create a new file
		f = excelize.NewFile()
	} else {
		// Some other error occurred when checking file existence
		log.Printf("ERROR: Failed to check Excel file existence: %s\n", err)
		return nil, err
	}
	return f, nil
}

// AddSheetToExcelFile adds a new sheet to the Excel file with the given headers.
func AddSheetToExcelFile(f *excelize.File, sheetName string, headers []string) error {

	// Check if the sheet already exists
	idx, err := f.GetSheetIndex(sheetName)
	if err != nil {
		log.Printf("ERROR: Cloud not check if sheet already exist: %s\n", err)
		return err
	} else if idx != -1 {
		log.Println("Sheet already exists, no need to add it again")
		return nil
	}

	// Add a new sheet
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Printf("ERROR: Cloud not create sheet: %s\n", err)
		return err
	}
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1) // 1-based indexing
		if err != nil {
			log.Printf("ERROR: Failed to convert coordinates to cell name: %s\n", err)
			return err
		}
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			log.Printf("ERROR: Failed to set cell value: %s\n", err)
			return err
		}
	}
	f.SetActiveSheet(index)
	log.Printf("INFO: Added new sheet '%s' to Excel file\n", sheetName)

	err = f.DeleteSheet("Sheet1")

	if err != nil {
		log.Printf("ERROR: Cloud not delete sheet1: %s\n", err)
		return err
	}

	return nil
}
