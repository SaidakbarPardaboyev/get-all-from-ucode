package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	get_all "github.com/SaidakbarPardaboyev/get-all-from-ucode"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		startingTime = time.Now()
		endingTime   time.Time
	)

	// Read database credentials from environment variables
	cfg := get_all.Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_TYPE:     os.Getenv("DB_TYPE"),
	}

	apis, err := get_all.New(&cfg)
	if err != nil {
		panic(fmt.Errorf("error creating function: %v", err))
	}

	items, err := apis.Items("products").GetAll().Pipeline([]map[string]any{
		{
			"$match": map[string]any{
				"guid": []string{"340845d9-0219-4aec-aae1-f28c5c4bb859", "3fc80387-01aa-4c28-bf51-d71bf7ba8a4e"},
			},
		},
		{
			"$lookup": map[string]any{
				"from":         "product_images",
				"localField":   "guid",
				"foreignField": "product_id",
				"as":           "product_images",
			},
		},
		{
			"$lookup": map[string]any{
				"from":         "product_options",
				"localField":   "guid",
				"foreignField": "product_id",
				"as":           "product_options",
			},
		},
		{
			"$lookup": map[string]any{
				"from":         "product_variations",
				"localField":   "guid",
				"foreignField": "product_id",
				"as":           "product_variations",
			},
		},
	}).Exec()
	if err != nil {
		panic(fmt.Errorf("error executing function: %v", err))
	}

	for _, item := range items {
		decidedData, _ := json.Marshal(item)
		fmt.Println(string(decidedData))
	}

	endingTime = time.Now()
	fmt.Println("Execution time:", endingTime.Sub(startingTime))
}
