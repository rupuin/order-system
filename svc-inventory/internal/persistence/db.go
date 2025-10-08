package persistence

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository() *Repository {
	return &Repository{
		db: initDb(),
	}
}

func (r *Repository) FetchItemStatus(itemID string) (string, error) {
	query := `
	SELECT status
	FROM items
	WHERE id = $1
	`

	var status string
	err := r.db.QueryRow(context.Background(), query, itemID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func initDb() *pgx.Conn {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("Missing database URL from env")
	}

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
