package main

import (
     "daily_check_in/excel"
	"fmt"
	"daily_check_in/api"
	"log"
	"github.com/joho/godotenv"
	"daily_check_in/api"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loading config...")	
	cfg, err := api.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Printf("Error creating API client: %v\n", err)
		return
	}
	data := excel.FormatToolBoxTalkData(client)
	fmt.Printf("Formatted Toolbox Talk Data: %+v\n", data)

	fmt.Println("Creating Crew Allocation Sheet...")
	excel.CreateCrewAllocationSheet("Crew_Sheet.xlsx")
	fmt.Println("Crew_Sheet.xlsx created successfully")
}