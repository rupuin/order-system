package persistence

import (
	"context"
	"fmt"
	"log"
	"os"
	"svc-order/dto"

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

func (r *Repository) CreateOrder(itemID, address, status string) (*dto.Order, error) {
	query := `
	INSERT INTO orders (item_id, address, status)
	VALUES ($1, $2, $3)
	RETURNING (id, item_id, address, status, created_at, updated_at)
	`
	var order dto.Order
	err := r.db.QueryRow(context.Background(), query, itemID, address, status).Scan(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *Repository) UpdateOrderStatus(orderID, newStatus string) (int, error) {
	query := `
	UPDATE orders
	WHERE id = ($1)
	SET status = ($2)
	RETURNING id
	`

	var updatedOrderId int
	err := r.db.QueryRow(context.Background(), query, orderID, newStatus).Scan(&updatedOrderId)
	if err != nil {
		return 0, err
	}
	return updatedOrderId, nil
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
