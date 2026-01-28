package main

import (
	"daily_check_in/api"
	"daily_check_in/excel"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Loading config...")
	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetching toolbox talks:")

	allCrews, err := client.GetCrewAllocationData()
	if err != nil {
		log.Fatal(err)
	}

	if allCrews == nil {
		log.Fatal("No crew allocation data retrieved")
	}
	fmt.Println("amount of crews fetched:", len(allCrews))
	jsonData, _ := json.MarshalIndent(allCrews, "", "  ")
	fmt.Println(string(jsonData))
	fmt.Println("Creating Excel file...")
	if err := excel.CreateCrewAllocationSheet("Crew_Allocation_Recap.xlsx", allCrews); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Excel file created: Crew_Allocation_Recap.xlsx")
}
