package main

import (
	"prg_tools/material"
)

func main() {
	svc := material.NewTestProgressEstimateService()

	fromDate := "2024-01-01"
	toDate := "2024-05-31"
	materialLog, err := svc.GetMaterialLogs(fromDate, toDate)
	if err != nil {
		panic(err)
	}
	if err := svc.ExportMaterialLogs(materialLog); err != nil {
		panic(err)
	}
}
