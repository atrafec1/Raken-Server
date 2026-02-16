package payroll

import (
	"prg_tools/payroll/adapter/cp"
	"prg_tools/payroll/adapter/excel"
	"prg_tools/payroll/adapter/raken"
)

func NewCPService() (*PayrollService, error) {
	reader, err := raken.NewRakenAPIAdapter()
	if err != nil {
		return nil, err
	}
	excelExporter := excel.NewPayrollExcelExporter("I:\\Daily Time Sheets\\Raken Timesheets", "I:\\Daily Time Sheets\\Raken Timesheets\\Warnings")
	payrollExporter := cp.NewAdapter("T:\\CP\\CPData")

	return NewPayrollService(reader, payrollExporter, excelExporter), nil
}

func NewTestCPService() (*PayrollService, error) {
	testDir := "./test_output"
	reader, err := raken.NewRakenAPIAdapter()
	if err != nil {
		return nil, err
	}

	excelExporter := excel.NewPayrollExcelExporter(testDir, testDir)
	payrollExporter := cp.NewAdapter(testDir)
	return NewPayrollService(reader, payrollExporter, excelExporter), nil
}
