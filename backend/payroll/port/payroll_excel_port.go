package port

import "prg_tools/payroll/dto"

//outbound port for exporting payroll entries to excel
type PayrollExcelPort interface {
	ExportWarnings([]dto.Warning) error
	ExportPayrollEntries([]dto.PayrollEntry) error
}
