package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect() (*pgx.Conn, error) {
	link := "postgres://a4bhi:a4bhi@localhost:5432/vaultx"

	conn, err := pgx.Connect(context.Background(), link)

	return conn, err
}
