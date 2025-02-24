package storage

import (
	"context"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg/db"
	"github.com/SaidakbarPardaboyev/get-all-from-ucode/storage/repo"

	"github.com/SaidakbarPardaboyev/get-all-from-ucode/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaidakbarApis interface {
	Config() *pkg.Config

	Items(collection string) repo.ItemsI
}

func New(cfg *pkg.Config) (SaidakbarApis, error) {
	var (
		mongoDatabase *mongo.Database
		conn          *mongo.Client
		postgresConn  *pgxpool.Pool
		err           error
	)

	if cfg.DB_TYPE == "mongo" {
		conn, mongoDatabase, err = db.ConnectMongoDB(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		postgresConn, err = db.ConnectPostgresDB(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &object{
		config:      cfg,
		mongoDb:     mongoDatabase,
		mongoClient: conn,
		postgresDb:  postgresConn,
	}, nil
}

type object struct {
	config      *pkg.Config
	mongoDb     *mongo.Database
	mongoClient *mongo.Client
	postgresDb  *pgxpool.Pool
}

func (o *object) Config() *pkg.Config {
	return o.config
}

func (o *object) Items(collection string) repo.ItemsI {
	return &repo.APIItem{
		Collection: collection,
		Config: &pkg.InnerConfig{
			MongoDb:    o.mongoDb,
			PostgresDb: o.postgresDb,
			DB_TYPE:    o.config.DB_TYPE,
		},
	}
}

func (o *object) Close(ctx context.Context) error {
	return o.mongoClient.Disconnect(ctx)
}
