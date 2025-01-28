package get_all

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

type innerConfig struct {
	MongoDb    *mongo.Database
	PostgresDb *pgxpool.Pool

	// mongo or postgress
	DB_TYPE string
}

type APIItem struct {
	collection string
	config     *innerConfig
}

type GetAllI struct {
	collection string
	config     *innerConfig
	sort       map[string]interface{}
	filter     map[string]interface{}
	limit      int64
	skip       int64
}

type object struct {
	config     *Config
	mongoDb    *mongo.Database
	postgresDb *pgxpool.Pool
}
