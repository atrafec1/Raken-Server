package main

import (
	"fmt"
	"prg_tools/payroll"
)

func main() {
	cpPayrollService, err := payroll.NewTestCPService()

	if err != nil {
		panic(err)
	}

	payrollEntryResult, err := cpPayrollService.GetEntries("2026-02-09", "2026-02-14")
	if err != nil {
		panic(err)
	}
	fmt.Println("Payroll entries retrieved successfully")

	if err := cpPayrollService.ExportExcel(payrollEntryResult); err != nil {
		panic(err)
	}

	fmt.Println("Payroll entries exported to excel successfully")

	if err := cpPayrollService.ExportToPayroll(payrollEntryResult.Entries); err != nil {
		panic(err)
	}

	fmt.Println("Payroll entries exported to payroll system successfully")
}
