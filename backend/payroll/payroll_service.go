package payroll

import (
	"fmt"
	"prg_tools/payroll/dto"
	"prg_tools/payroll/port"
	"errors"
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

func (s *PayrollService) GetEntries(fromDate, toDate string) (port.PayrollEntryResult, error) {
	result, err := s.EntryReader.GetPayrollEntries(fromDate, toDate)
	if err != nil {
		return port.PayrollEntryResult{}, fmt.Errorf("failed to get payroll entries: %w", err)
	}
	return result, nil
}


