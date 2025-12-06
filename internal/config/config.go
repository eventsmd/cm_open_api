package config

import (
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	PostgresConnStr string
	DB              *pgxpool.Pool
}

func LoadConfig() (*Config, error) {

	return &Config{
		PostgresConnStr: os.Getenv("POSTGRES_CONN_STR"),
	}, nil
}
