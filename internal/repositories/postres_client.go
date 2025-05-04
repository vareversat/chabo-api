package repositories

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var postgresDSN = os.Getenv("POSTGRES_DSN")

func NewPostgresClient() *pgxpool.Pool {
	var ctx = context.Background()
	pool, err := pgxpool.New(context.Background(), postgresDSN)
	if err != nil {
		panic(err)
	} else {
		// Verify the connection
		if err := pool.Ping(ctx); err != nil {
			panic(err)
		}
	}

	return pool
}
