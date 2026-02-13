package main

import (
	"fmt"

	"prg_tools/payroll"
	"prg_tools/payroll/adapter/excel"
	"prg_tools/payroll/adapter/raken"
)

func main() {
	fmt.Println("starting")
	raken_adapter, err := raken.NewRakenAPIAdapter()
	if err != nil {
		panic(err)
	}
	excel_exporter := excel.NewExcelPayrollExporter("Desktop")
	payroll_service := payroll.NewPayrollService(raken_adapter, excel_exporter)
	entries, err := payroll_service.GetEntries("2026-01-26", "2026-02-01")
	if err != nil {
		panic(err)
	}
	fmt.Printf("warnings: %+v\n", entries.Warnings)
	if err := payroll_service.Export(entries.Entries); err != nil {
		panic(err)
	}
	if err := payroll_service.ExportWarnings(entries.Warnings, ""); err != nil {
		panic(err)
	}
}
