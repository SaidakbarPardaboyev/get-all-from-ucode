package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"github.com/SaidakbarPardaboyev/get-all-from-ucode/storage"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	GUID               string             `bson:"guid"`
	ShopifyID          int64              `bson:"shopify_id"`
	Title              string             `bson:"title"`
	Handle             string             `bson:"handle"`
	BodyHTML           string             `bson:"body_html"`
	ProductType        string             `bson:"product_type"`
	PublishedAt        string             `bson:"published_at"`
	PublishedScope     []string           `bson:"published_scope"`
	Status             []string           `bson:"status"`
	Tags               string             `bson:"tags"`
	Vendor             string             `bson:"vendor"`
	PlatformID         string             `bson:"platform_id"`
	PlatformType       []string           `bson:"platform_type"`
	Disabled           bool               `bson:"disabled"`
	DisabledFrom       []string           `bson:"disabled_from"`
	CreatedTime        string             `bson:"created_time"`
	UpdatedTime        string             `bson:"updated_time"`
	MerchantID         string             `bson:"merchant_id"`
	PlatformCategoryID string             `bson:"platform_category_id"`
	ShopifyProductID   *int64             `bson:"shopify_product_id,omitempty"`
	ShopifyVariantID   *int64             `bson:"shopify_variant_id,omitempty"`
	CreatedAt          *time.Time         `bson:"createdAt"`
	UpdatedAt          *time.Time         `bson:"updatedAt"`
	Version            int                `bson:"__v"`
}

type Merchant struct {
	ID                                primitive.ObjectID `bson:"_id,omitempty"`
	GUID                              string             `bson:"guid"`
	Logo                              string             `bson:"logo"`
	XAPIKey                           string             `bson:"x_api_key"`
	IsActive                          bool               `bson:"is_active"`
	Interval                          int32              `bson:"interval,omitempty"`
	CronjobLastPerformTime            string             `bson:"cronjob_last_perform_time"`
	ResponseBodyExample               string             `bson:"response_body_example"`
	MyBazarAfterCreateMerchantDisable bool               `bson:"mybazar-after-create-merchant_disable"`
	SyncType                          []string           `bson:"sync_type"`
	IntegrationType                   []string           `bson:"integration_type"`
	Name                              string             `bson:"name"`
	WarehouseLocationID               int64              `bson:"warehouse_location_id"`
	DisableAllProducts                bool               `bson:"disable_all_products"`
	// CreatedAt                         *time.Time         `bson:"createdAt"`
	// UpdatedAt                         *time.Time         `bson:"updatedAt"`
	Version int `bson:"__v"`
}

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

	var products = []Merchant{}
	err = apis.Items("merchants").GetAll().Pipeline([]map[string]any{}).Exec(&products)
	if err != nil {
		panic(fmt.Errorf("error getting products: %v", err))
	}

	decodedData, _ := json.Marshal(products)
	fmt.Println(string(decodedData))

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
