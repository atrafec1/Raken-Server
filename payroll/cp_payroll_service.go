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
		fmt.Printf("Employee: %s, Date: %s, Job Number: %s, Cost Code: %s, Pay Type: %s, Hours: %.2f\n", entry.EmployeeCode,
			entry.CurrentDate, entry.JobNumber, entry.CostCode)
	}
	return nil
}
