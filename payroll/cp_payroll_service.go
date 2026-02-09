package payroll

import (
	"daily_check_in/payroll/port"
	"fmt"
)

type CPPayrollService struct {
	PayrollEntryPort port.PayrollEntryPort
}

func NewCPPayrollService(payrollEntryPort port.PayrollEntryPort) *CPPayrollService {
	return &CPPayrollService{
		PayrollEntryPort: payrollEntryPort,
	}
}

func (s *CPPayrollService) PrintPayrollEntries(fromDate, toDate string) error {

	entries, err := s.PayrollEntryPort.GetPayrollEntries(fromDate, toDate)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		fmt.Printf("Date: %s, Employee code: %s, Regular hours: %.2f, Premium Hours: %.2f, Overtime Hours %.2f \n", entry.CurrentDate, entry.EmployeeCode, entry.RegularHours, entry.PremiumHours, entry.OvertimeHours)
	}
	return nil
}
