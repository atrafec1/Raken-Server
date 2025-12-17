package excel

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// MyExcel wraps an excelize.File
type MyExcel struct {
	File *excelize.File
}

// SetHeaderValues fills Sheet1 with predefined headers and today's date
func (m *MyExcel) SetHeaderValues() error {
	currentTime := time.Now()
	year, month, day := currentTime.Date()
	todayStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	// Map of cell -> value
	cellValues := map[string]string{
		"B1": "TIMESHEET REVIEW / RECAP",
		"F1": todayStr,
		"J2": "",
		"K2": "",
		"B5": "",
		"C5": "",
		"D5": "Last Name",
		"E5": "First Name",
		"F5": "Class",
		"G5": "Equip.",
		"H5": "",
		"I5": "",
		"J5": "",
		"K5": "",
		"F4": "M",
		"G4": "T",
		"H4": "W",
		"I4": "Th",
		"J4": "F",
		"K4": "S",
	}

	for cell, value := range cellValues {
		if err := m.File.SetCellValue("Sheet1", cell, value); err != nil {
			return fmt.Errorf("failed to set cell %s: %v", cell, err)
		}
	}

	return nil
}

// CreateCrewAllocationSheet creates a new workbook, sets headers, and saves it
func CreateCrewAllocationSheet(filename string) error {
	f := excelize.NewFile()

	// Wrap with MyExcel
	myExcel := &MyExcel{File: f}

	// Set header values
	if err := myExcel.SetHeaderValues(); err != nil {
		return err
	}

	// Example: create a second sheet
	sheet2 := "Sheet2"
	f.NewSheet(sheet2)
	f.SetCellValue(sheet2, "A2", "This is on Sheet 2!")

	// Save the workbook
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}

	// Close file
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}

	return nil
}

