package main

import (
	"daily_check_in/payroll"
	"daily_check_in/payroll/adapter"
)

func main() {
	adapter, err := adapter.NewRakenAPIAdapter()
	if err != nil {
		panic(err)
	}
	service := payroll.NewCPPayrollService(adapter)
	if err := service.PrintPayrollEntries("2026-01-01", "2026-01-31"); err != nil {
		panic(err)
	}

}
