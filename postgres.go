package get_all

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgresDB(cfg *Config) (*pgxpool.Pool, error) {
	// Build the PostgreSQL connection string
	postgresString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME,
	)

	// Create a connection pool
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, postgresString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL: %v", err)
	}

	// Test the connection
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging PostgreSQL: %v", err)
	}

	fmt.Println("Successfully connected and pinged PostgreSQL!")

	return conn, nil
}
