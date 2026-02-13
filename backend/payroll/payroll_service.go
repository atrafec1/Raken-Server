package payroll

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"prg_tools/payroll/dto"
	"prg_tools/payroll/port"

	"github.com/xuri/excelize/v2"
)

type PayrollService struct {
	EntryReader port.PayrollEntryPort
	Exporter    port.PayrollExportPort
}

func NewPayrollService(payrollEntryPort port.PayrollEntryPort, payrollExportPort port.PayrollExportPort) *PayrollService {
	return &PayrollService{
		EntryReader: payrollEntryPort,
		Exporter:    payrollExportPort,
	}
}

func (s *PayrollService) Export(payrollEntries []dto.PayrollEntry) error {
	if len(payrollEntries) == 0 {
		return errors.New("no payroll entries")
	}

	if err := s.Exporter.ExportPayrollEntries(payrollEntries); err != nil {
		return fmt.Errorf("failed to export payroll entries: %w", err)
	}
	return nil
}

func (s *PayrollService) GetEntries(fromDate, toDate string) (dto.PayrollEntryResult, error) {
	result, err := s.EntryReader.GetPayrollEntries(fromDate, toDate)
	if err != nil {
		return dto.PayrollEntryResult{}, fmt.Errorf("failed to get payroll entries: %w", err)
	}
	return result, nil
}

func (s *PayrollService) ExportWarnings(warnings []dto.Warning, path string) error {
	if len(warnings) == 0 {
		return fmt.Errorf("no warnings to export")
	}
	var savePath string
	if path == "" {
		//default to downloads
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		savePath = filepath.Join(homeDir, "Downloads", fmt.Sprintf("payroll_warnings_.xlsx"))
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Warnings"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return fmt.Errorf("failed to create sheet: %w", err)
	}

	// Set headers
	f.SetCellValue(sheetName, "A1", "Warning Type")
	f.SetCellValue(sheetName, "B1", "Message")

	// Write warning data
	for i, warning := range warnings {
		row := i + 2 // Start at row 2 (after header)
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), warning.WarningType)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), warning.Message)
	}

	// Set active sheet and delete default Sheet1
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Save file
	if err := f.SaveAs(savePath); err != nil {
		return fmt.Errorf("failed to save excel file: %w", err)
	}

	return nil
}
