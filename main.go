package main

import (
	// "daily_check_in/api"
	// "log"
	"fmt"
	// "github.com/joho/godotenv"
	"daily_check_in/excel"
)

// func main() {
// 	if err := godotenv.Load(".env"); err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	cfg, err := api.LoadConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Config loaded successfully:", cfg)
// 	client, err := api.NewClient(cfg)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fromDate := "2025-12-01"
// 	toDate := "2025-12-07"
// 	fmt.Println("Fetching timecards from", fromDate, "to", toDate)
//     timecardsResp, err := client.GetTimecards(fromDate, toDate)	
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Timecards Retrieved:")
// 	fmt.Println(timecardsResp)
// 	fmt.Println(len(timecardsResp.Collection))	
// }

func main() {
	fmt.Println("Creating excel sheet!")
	if err := excel.CreateCrewAllocationSheet("fake_excel.xlsx"); err != nil {
		fmt.Printf("Error creating excel sheet: %v", err)
		return
	}
	fmt.Println("Success!")
}
