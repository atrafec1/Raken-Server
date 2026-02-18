package main

import (
	"fmt"
	"prg_tools/material"
)

func main() {
	svc := material.RakenProgressEstimateService("./test_output/raken_estimate")

	fromDate := "2026-01-02"
	toDate := "2026-01-31"
	materialLog, err := svc.GetMaterialLogs(fromDate, toDate)
	if err != nil {
		panic(err)
	}
	fmt.Println("Material logs retrieved successfully:")
	if err := svc.ExportMaterialLogs(materialLog); err != nil {
		panic(err)
	}
	fmt.Println("Material logs exported successfully")
}
