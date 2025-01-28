package get_all

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type SaidakbarApis interface {
	Config() *Config

	Items(collection string) ItemsI
}

func New(cfg *Config) (SaidakbarApis, error) {
	var (
		mongoDatabase *mongo.Database
		postgresConn  *pgxpool.Pool
		err           error
	)

	if cfg.DB_TYPE == "mongo" {
		mongoDatabase, err = ConnectMongoDB(cfg)
		if err != nil {
			return nil, err
		}
	} else {
		postgresConn, err = ConnectPostgresDB(cfg)
		if err != nil {
			return nil, err
		}
	}

	return &object{
		config:     cfg,
		mongoDb:    mongoDatabase,
		postgresDb: postgresConn,
	}, nil
}

func (o *object) Config() *Config {
	return o.config
}
