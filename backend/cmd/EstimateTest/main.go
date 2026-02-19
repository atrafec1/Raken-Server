package main

import (
	"prg_tools/material"
)

func main() {
	svc := material.NewTestProgressEstimateService()

	fromDate := "2024-01-01"
	toDate := "2024-05-31"
	matInfo, err := svc.GetJobMaterialInfo(fromDate, toDate)
	if err != nil {
		panic(err)
	}
	if err := svc.ExportJobMaterialInfo(matInfo); err != nil {
		panic(err)
	}
}
