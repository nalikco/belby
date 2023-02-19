package entities

import "time"

type Product struct {
	ID        int64
	MessageID int64
	Shop      string
	Title     string
	Price     float64
	Link      string
	CreatedAt time.Time
}
