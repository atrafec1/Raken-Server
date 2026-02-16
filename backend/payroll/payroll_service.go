package payroll

import (
	"errors"
	"fmt"
	"prg_tools/payroll/dto"
	"prg_tools/payroll/port"
)

type PayrollService struct {
	entryReader     port.PayrollEntryPort
	payrollExporter port.PayrollExportPort
	excelExporter   port.PayrollExcelPort
}

func NewPayrollService(
	payrollEntryPort port.PayrollEntryPort,
	payrollExportPort port.PayrollExportPort,
	excelExporter port.PayrollExcelPort) *PayrollService {

	return &PayrollService{
		entryReader:     payrollEntryPort,
		payrollExporter: payrollExportPort,
		excelExporter:   excelExporter,
	}
}

func (s *PayrollService) ExportToPayroll(payrollEntries []dto.PayrollEntry) error {
	if len(payrollEntries) == 0 {
		return errors.New("no payroll entries")
	}

	if err := s.payrollExporter.ExportPayrollEntries(payrollEntries); err != nil {
		return fmt.Errorf("failed to export payroll entries: %w", err)
	}
	return nil
}

func (s *PayrollService) ExportWarnings(warnings []dto.Warning) error {
	if len(warnings) == 0 {
		return fmt.Errorf("no warnings to export")
	}
	if err := s.excelExporter.ExportWarnings(warnings); err != nil {
		return fmt.Errorf("failed to export warnings: %w", err)
	}
	return nil
}

func (s *PayrollService) ExportToExcel(payrollEntries []dto.PayrollEntry) error {
	if len(payrollEntries) == 0 {
		return errors.New("no payroll entries")
	}
	if err := s.excelExporter.ExportPayrollEntries(payrollEntries); err != nil {
		return fmt.Errorf("failed to export payroll entries to excel: %w", err)
	}
	return nil
}

// ExportExcel exports both payroll entries and warnings to excel files. Warnings are exported only if there are any.
func (s *PayrollService) ExportExcel(result dto.PayrollEntryResult) error {
	if err := s.excelExporter.ExportPayrollEntries(result.Entries); err != nil {
		return fmt.Errorf("failed to export payroll entries to excel: %w", err)
	}
	if err := s.excelExporter.ExportWarnings(result.Warnings); err != nil {
		return fmt.Errorf("failed to export warnings to excel: %w", err)
	}
	return nil
}

func (s *PayrollService) GetEntries(fromDate, toDate string) (dto.PayrollEntryResult, error) {
	result, err := s.entryReader.GetPayrollEntries(fromDate, toDate)
	if err != nil {
		return dto.PayrollEntryResult{}, fmt.Errorf("failed to get payroll entries: %w", err)
	}
	return result, nil
}
