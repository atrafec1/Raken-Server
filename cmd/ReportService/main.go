package main

import (
	reports "daily_check_in/Reports"
	"daily_check_in/api"
	"log"
)

func main() {
	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	apiClient, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	filePath, err := reports.ExportReports("2026-01-01", "2026-01-31", ".", apiClient)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved report file: %s\n", filePath)

}
