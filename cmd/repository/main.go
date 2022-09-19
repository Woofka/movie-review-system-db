package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	reviewPkg "gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/core/review"
	"gitlab.ozon.dev/Woofka/movie-review-system-db/internal/pkg/repository/postgres"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	psqlConn := "host=localhost port=5432 user=user password=password dbname=movie_review sslmode=disable"
	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	pgConfig := pool.Config()
	pgConfig.MaxConnIdleTime = time.Minute
	pgConfig.MaxConnLifetime = time.Hour
	pgConfig.MinConns = 2
	pgConfig.MaxConns = 4

	var review reviewPkg.Interface
	review = reviewPkg.New(postgres.New(pool))

	err = runGRPCServer(review)
	if err != nil {
		log.Fatal(err)
	}
}
