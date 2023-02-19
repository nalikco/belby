package storages

import (
	"belby/internal/entities"
	"belby/pkg/database"
	"time"
)

type ProductsStorage struct {
	db *database.Database
}

func NewProductsStorage(db *database.Database) *ProductsStorage {
	return &ProductsStorage{
		db: db,
	}
}

func (s *ProductsStorage) Create(product *entities.Product) error {
	return s.db.Conn.
		QueryRow(s.db.Ctx, "INSERT INTO products(message_id, shop, title, price, link, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
			product.MessageID, product.Shop, product.Title, product.Price, product.Link, time.Now()).
		Scan(&product.ID)
}
