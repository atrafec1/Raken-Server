package main

import (
     "daily_check_in/excel"
	"fmt"
)

func main() {
	fmt.Println("Creating Crew Allocation Sheet...")
	excel.CreateCrewAllocationSheet("Crew_Sheet.xlsx")
	fmt.Println("Crew_Sheet.xlsx created successfully")
}