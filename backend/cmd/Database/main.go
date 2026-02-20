package main

import (
	"fmt"
	"log"
	"prg_tools/database"
	"prg_tools/external/rakenapi/db"
)

func main() {
	fmt.Println("intitializing db")
	database.InitTestDB()
	fmt.Println("Running data ingestion")
	pipeline, err := db.NewPipeline()
	if err != nil {
		log.Fatalf("Error initializing pipeline: %v", err)
	}
	pipeline.LoadDatabase()
	fmt.Println("Data ingestion complete. Printing schema and data:")
	printSchema()
	printData()
}
func printData() {
	var jobs []database.Job
	if err := database.DB.Find(&jobs).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("Jobs in DB:", jobs)

	var employees []database.Employee
	if err := database.DB.Find(&employees).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("Employees in DB:", employees)

	var costCodes []database.CostCode
	if err := database.DB.Find(&costCodes).Error; err != nil {
		log.Fatal(err)
	}
	fmt.Println("CostCodes in DB:", costCodes)
}

func printSchema() {
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := sqlDB.Query("SELECT name, type, sql FROM sqlite_master WHERE type='table';")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, typ, sqlStmt string
		if err := rows.Scan(&name, &typ, &sqlStmt); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Table: %s\nType: %s\nSQL: %s\n\n", name, typ, sqlStmt)
	}
}
