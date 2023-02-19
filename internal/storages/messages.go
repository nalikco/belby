package storages

import (
	"belby/internal/entities"
	"belby/pkg/database"
	"time"
)

type MessagesStorage struct {
	db *database.Database
}

func NewMessagesStorage(db *database.Database) *MessagesStorage {
	return &MessagesStorage{
		db: db,
	}
}

func (s *MessagesStorage) Create(message *entities.Message) error {
	return s.db.Conn.
		QueryRow(s.db.Ctx, "INSERT INTO messages(vk_id, user_id, message, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
			message.VkID, message.UserID, message.Message, time.Now()).
		Scan(&message.ID)
}
