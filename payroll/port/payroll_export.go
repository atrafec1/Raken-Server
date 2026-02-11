package port

import "daily_check_in/payroll/dto"

type PayrollExportPort interface {
	ExportPayrollEntries([]dto.PayrollEntry) error
}
