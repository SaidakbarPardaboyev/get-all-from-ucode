package pkg

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	// mongo or postgress
	DB_TYPE string
}

type InnerConfig struct {
	MongoDb    *mongo.Database
	PostgresDb *pgxpool.Pool

	// mongo or postgress
	DB_TYPE string
}

type APIItem struct {
	Collection string
	Config     *InnerConfig
}

type GetAllI struct {
	Collection string
	Config     *InnerConfig
	Sort       map[string]interface{}
	Filter     map[string]interface{}
	Limit      int64
	Skip       int64
	Pipeline   []map[string]any
}
