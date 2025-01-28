package main

import (
	"fmt"
	"get_all"
	"time"
)

func main() {
	var (
		startingTime = time.Now()
		endingTime   time.Time
	)

	// Call the handler directly
	var (
		// cfg = get_all.Config{
		// 	DB_HOST:     "142.93.164.37",
		// 	DB_PORT:     "27017",
		// 	DB_USER:     "mybazar_0f1dc8ebeda749bb80a31f78af44536d_p_obj_build_svcs",
		// 	DB_PASSWORD: "JurgVHHYW5",
		// 	DB_NAME:     "mybazar_0f1dc8ebeda749bb80a31f78af44536d_p_obj_build_svcs",
		// 	DB_TYPE:     "mongo",
		// }

		cfg = get_all.Config{
			DB_HOST:     "142.93.164.37",
			DB_PORT:     "30032",
			DB_USER:     "brrauf_3a40b209092f45eca6a93a8d8f1af9d4_p_postgres_svcs",
			DB_PASSWORD: "BnwgaR4Hvk",
			DB_NAME:     "brrauf_3a40b209092f45eca6a93a8d8f1af9d4_p_postgres_svcs",
			DB_TYPE:     "postgres",
		}
	)
	apis, err := get_all.New(&cfg)
	if err != nil {
		panic(fmt.Errorf("error creating function: %v", err))
	}

	items, err := apis.Items("request").GetAll().Count()
	if err != nil {
		panic(fmt.Errorf("error executing function: %v", err))
	}

	// fmt.Println(len(items))
	// for _, item := range items {
	// 	fmt.Println(item["id"], item["created_at"])
	// }
	fmt.Println(items)

	endingTime = time.Now()
	fmt.Println("Execution time:", endingTime.Sub(startingTime))
}
