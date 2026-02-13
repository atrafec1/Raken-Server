package port

import "prg_tools/payroll/dto"

type PayrollEntryPort interface {
	GetPayrollEntries(fromDate, toDate string) (dto.PayrollEntryResult, error)
}
