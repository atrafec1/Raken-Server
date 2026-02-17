package main

import (
	"fmt"
	"prg_tools/material"
	"prg_tools/material/domain"
)

func main() {
	progSvc := material.NewTestProgressEstimateService()

	logs, err := progSvc.GetMaterialLogs("2026-02-09", "2026-02-14")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved %d material log collections\n", len(logs))
	if err := progSvc.ExportMaterialLogs(logs); err != nil {
		panic(err)
	}
	fmt.Println("Material logs exported successfully")
}

func printMaterialLogs(logs []domain.MaterialLogCollection) {
	for _, log := range logs {
		fmt.Printf("Job: %s, From: %s, To: %s\n", log.Job.Name, log.FromDate, log.ToDate)
		for _, entry := range log.Logs {
			fmt.Printf("  Date: %s, Material: %s, Quantity: %.2f %s\n",
				entry.Date, entry.Material.Name, entry.Quantity, entry.Material.Unit)
		}
	}
}
