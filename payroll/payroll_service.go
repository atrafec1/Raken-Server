package payroll

import (
	"daily_check_in/payroll/dto"
	"daily_check_in/payroll/port"
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
		return errors.New("failed to export payroll entries")
	}
	return nil
}

func (s *PayrollService) GetEntries(fromDate, toDate string) (port.PayrollEntryResult, error) {
	result, err := s.EntryReader.GetPayrollEntries(fromDate, toDate)
	if err != nil {
		return port.PayrollEntryResult{}, errors.New("failed to get payroll entries")
	}
	return result, nil
}
