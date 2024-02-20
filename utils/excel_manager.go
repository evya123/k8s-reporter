// utils/excel_manager.go

package utils

import (
	"sync"

	"log"

	"github.com/xuri/excelize/v2"
)

// excelFileManager manages a singleton instance of an Excel file.
type excelFileManager struct {
	file *excelize.File
	once sync.Once
}

var manager *excelFileManager
var once sync.Once

// GetExcelFileManager returns the singleton instance of the Excel file manager.
func GetExcelFileManager() *excelFileManager {
	once.Do(func() {
		manager = &excelFileManager{}
	})
	return manager
}

// OpenOrCreateExcelFile checks if an Excel file exists and opens it, or creates a new file if it doesn't exist.
func (m *excelFileManager) OpenOrCreateExcelFile(filePath string) error {
	var err error
	m.once.Do(func() {
		m.file, err = OpenOrCreateExcelFile(filePath)
	})
	if err != nil {
		log.Printf("ERROR: Failed to open or create Excel file: %s\n", err)
		return err
	}
	return nil
}

// GetExcelFile returns the Excel file instance.
func (m *excelFileManager) GetExcelFile() *excelize.File {
	return m.file
}

// SaveExcelFile saves the Excel file to the provided path.
func (m *excelFileManager) SaveExcelFile(filePath string) error {
	if m.file != nil {
		if err := m.file.SaveAs(filePath); err != nil {
			log.Printf("ERROR: Failed to save the Excel file: %s\n", err)
			return err
		}
		log.Println("INFO: Excel file saved successfully.")
		return nil
	}
	log.Println("ERROR: No Excel file to save.")
	return nil
}
