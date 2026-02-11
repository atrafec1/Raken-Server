package port

import "daily_check_in/payroll/dto"

type PayrollEntryPort interface {
	GetPayrollEntries(fromDate, toDate string) (PayrollEntryResult, error)
}

type PayrollEntryResult struct {
	Entries  []dto.PayrollEntry
	Warnings []dto.Warning
}
