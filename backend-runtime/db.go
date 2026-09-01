package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// NewPool crea el pool de conexiones a Postgres, leyendo la misma
// configuración que usa el backend en Node (.env con DATABASE_URL o
// DB_HOST/DB_PORT/DB_NAME/DB_USER/DB_PASSWORD).
func NewPool(ctx context.Context) (*pgxpool.Pool, error) {
	// Silencioso si no encuentra .env (por ejemplo en producción, donde las
	// variables ya vienen inyectadas por el entorno).
	_ = godotenv.Load()

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el pool de conexiones: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo conectar a Postgres: %w", err)
	}

	return pool, nil
}
