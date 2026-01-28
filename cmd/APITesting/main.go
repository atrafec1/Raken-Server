package main

import (
	"daily_check_in/api"
	"fmt"
	"log"

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

	fmt.Println("Fetching toolbox talks:")

	toolboxTalks, err := client.GetToolboxTalks()
	if err != nil {
		log.Fatal(err)
	}

	var employeeNames []string

	for _, talk := range toolboxTalks.Collection {

		for _, attendee := range talk.Attendees {
			employee, err := client.GetEmployeeByUUID(attendee.Member.UUID)
			if err != nil {
				log.Printf("Employee with UUID %s not found\n", attendee.Member.UUID)
			}
			fullName := fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
			employeeNames = append(employeeNames, fullName)
		}
	}
	fmt.Println("Attendees:")
	fmt.Println(employeeNames)
}
