package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"github.com/SaidakbarPardaboyev/get-all-from-ucode/storage"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read database credentials from environment variables
	cfg := pkg.Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_TYPE:     os.Getenv("DB_TYPE"),
	}

	apis, err := storage.New(&cfg)
	if err != nil {
		panic(fmt.Errorf("error creating function: %v", err))
	}

	err = apis.Items("products").MultipleUpdate().Exec(context.Background(), []mongo.WriteModel{mongo.NewUpdateOneModel().
		SetFilter(bson.M{
			"guid": bson.M{"$in": productGuids},
		}).
		SetUpdate(bson.M{
			"$set": bson.M{
				"shopify_id": float64(1),
				// "shopify_inventory_item_id": float64(variant.InventoryItemId),
			},
		})})
	if err != nil {
		panic("Error on update: " + err.Error())
	}

	// var (
	// 	startingTime = time.Now()
	// 	items        []map[string]any
	// 	uniqueGuids  = map[string]bool{}
	// )

	// // itemsResp, err := apis.Items("products").GetAll().Limit(50).Skip(int64(i) * 50).Pipeline([]map[string]any{
	// itemsResp, err := apis.Items("products").GetAll().Pipeline([]map[string]any{
	// 	{
	// 		"$match": map[string]any{
	// 			"guid": productGuids,
	// 			// "disabled": false,
	// 		},
	// 	},
	// 	{
	// 		"$lookup": map[string]any{
	// 			"from": "product_images",
	// 			"let":  map[string]any{"productId": "$guid"},
	// 			"pipeline": []map[string]any{
	// 				{
	// 					"$match": map[string]any{
	// 						"$expr": map[string]any{
	// 							"$and": []map[string]any{
	// 								{"$eq": []any{"$product_id", "$$productId"}},
	// 								{"$eq": []any{"$disabled", false}},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 			"as": "product_images",
	// 		},
	// 	},
	// 	{
	// 		"$lookup": map[string]any{
	// 			"from": "product_options",
	// 			"let":  map[string]any{"productId": "$guid"},
	// 			"pipeline": []map[string]any{
	// 				{
	// 					"$match": map[string]any{
	// 						"$expr": map[string]any{
	// 							"$and": []map[string]any{
	// 								{"$eq": []any{"$product_id", "$$productId"}},
	// 								{"$eq": []any{"$disabled", false}},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 			"as": "product_options",
	// 		},
	// 	},
	// 	{
	// 		"$lookup": map[string]any{
	// 			"from": "product_variations",
	// 			"let":  map[string]any{"productId": "$guid"},
	// 			"pipeline": []map[string]any{
	// 				{
	// 					"$match": map[string]any{
	// 						"$expr": map[string]any{
	// 							"$and": []map[string]any{
	// 								{"$eq": []any{"$product_id", "$$productId"}},
	// 								{"$eq": []any{"$disabled", false}},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			},
	// 			"as": "product_variations",
	// 		},
	// 	},
	// }).Exec()
	// if err != nil {
	// 	panic(fmt.Errorf("error executing function: %v", err))
	// }

	// items = append(items, itemsResp...)

	// fmt.Println("Execution time:", time.Since(startingTime))

	// // for _, item := range items {
	// // 	decidedData, _ := json.Marshal(item)
	// // 	fmt.Println(string(decidedData))
	// // }

	// for _, item := range items {
	// 	uniqueGuids[cast.ToString(item["guid"])] = true
	// }

	// fmt.Println("uniqueGuids", len(uniqueGuids))
}

var productGuids = []string{
	"8f0a3f2e-1674-46e1-a932-bbd7ddb68809",
}
