package main

import (
	"fmt"
	"prg_tools/material"
	"prg_tools/material/domain"
)

func main() {
	svc := material.NewTestProgressEstimateService()

	fromDate := "2024-01-01"
	toDate := "2024-05-31"
	materialLog, err := svc.GetMaterialLogs(fromDate, toDate)
	if err != nil {
		panic(err)
	}
	printMaterialLogs(materialLog)
	if err := svc.ExportMaterialLogs(materialLog); err != nil {
		panic(err)
	}
}

func printMaterialLogs(logs []domain.MaterialLogCollection) {
	for _, log := range logs {
		fmt.Printf("Job: %s, From: %s, To: %s\n", log.Job.Name, log.FromDate, log.ToDate)
		for _, entry := range log.Logs {
			fmt.Printf("Bid Item: %s,  Date: %s, Material: %s, Quantity: %.2f %s\n",
				entry.Material.BidNumber, entry.Date, entry.Material.Name, entry.Quantity, entry.Material.Unit)
		}
	}
}
