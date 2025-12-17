package main

import (
	"daily_check_in/api"
	"daily_check_in/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Get("/some-endpoint")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	log.Println("Status code:", resp.StatusCode)
}
