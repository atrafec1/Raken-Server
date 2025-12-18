package main

import (
	"daily_check_in/api"
	"log"
	"fmt"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Loading environment variables")
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loading config...")	
	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	fromDate := "2025-12-01"
	toDate := "2025-12-07"
	fmt.Println("Fetching timecards from", fromDate, "to", toDate)
    timecardsResp, err := client.GetTimecards(fromDate, toDate)	
	if err != nil {
		log.Fatal(err)
	}
	for _, timecard := range timecardsResp.Collection {
		employee, exists, err := client.GetEmployeeByUUID(timecard.Worker.UUID)
		if !exists {
			log.Fatalf("Employee with UUID %s not found", timecard.Worker.UUID)
		}else if err != nil {
			log.Fatalf("Error retrieving employee with UUID %s: %v", timecard.Worker.UUID, err)
		}
		fmt.Printf("Timecard for %s %s :\n", employee.FirstName, employee.LastName, )
	}
}

