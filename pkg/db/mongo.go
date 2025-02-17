package db

import (
	"context"
	"fmt"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(cfg *pkg.Config) (*mongo.Database, error) {
	mongoString := fmt.Sprintf("mongodb://%s:%s", cfg.DB_HOST, cfg.DB_PORT)

	credential := options.Credential{
		Username:   cfg.DB_USER,
		Password:   cfg.DB_PASSWORD,
		AuthSource: cfg.DB_NAME,
	}

	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoString).SetAuth(credential))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB %s", err.Error())
	}

	err = conn.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("error pinging MongoDB %s", err.Error())
	}

	fmt.Println("Successfully connected and pinged MongoDB")

	db := conn.Database(cfg.DB_NAME)

	return db, nil
}
