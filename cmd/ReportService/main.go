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

	reportService := reports.ReportService{
		Client: apiClient,
	}

	projects, err := reportService.GetProjectsWorkedOn("2026-01-01", "2026-01-31")
	if err != nil {
		log.Fatal(err)
	}
	allReports, err := reportService.GetReports("2026-01-01", "2026-01-31", projects)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Fetched reports for %d projects\n", len(allReports))

}
