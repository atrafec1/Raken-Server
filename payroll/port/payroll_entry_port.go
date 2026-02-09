package port

import "daily_check_in/payroll/dto"

type PayrollEntryPort interface {
	GetPayrollEntries(fromDate, toDate string) ([]dto.PayrollEntry, error)
}
