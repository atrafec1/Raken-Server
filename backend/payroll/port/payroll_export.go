package port

import "prg_tools/payroll/dto"

type PayrollExportPort interface {
	ExportPayrollEntries([]dto.PayrollEntry) error
	ExportPayrollWarnings([]dto.Warning) error
}
