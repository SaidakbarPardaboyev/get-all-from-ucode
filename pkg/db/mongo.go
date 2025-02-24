package db

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(cfg *pkg.Config) (*mongo.Client, *mongo.Database, error) {
	mongoString := fmt.Sprintf("mongodb://%s:%s", cfg.DB_HOST, cfg.DB_PORT)

	credential := options.Credential{
		Username:   cfg.DB_USER,
		Password:   cfg.DB_PASSWORD,
		AuthSource: cfg.DB_NAME,
	}

	clientOptions := options.Client().
		ApplyURI(mongoString).
		SetAuth(credential).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return conn, nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	if err = conn.Ping(ctx, nil); err != nil {
		return conn, nil, fmt.Errorf("error pinging MongoDB: %w", err)
	}

	fmt.Println("Successfully connected and pinged MongoDB")
	return conn, conn.Database(cfg.DB_NAME), nil
}
