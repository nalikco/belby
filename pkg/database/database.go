package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

type Database struct {
	config *Config
	Conn   *pgx.Conn
	Ctx    context.Context
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	SslMode  string
}

func NewDatabase(config *Config) (*Database, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return &Database{}, err
	}

	return &Database{
		config: config,
		Conn:   conn,
		Ctx:    ctx,
	}, nil
}

func (d *Database) Close() error {
	return d.Conn.Close(d.Ctx)
}
