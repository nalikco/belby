package storages

import (
	"belby/internal/entities"
	"belby/pkg/database"
	"time"
)

type UsersStorage struct {
	db *database.Database
}

func NewUsersStorage(db *database.Database) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (s *UsersStorage) GetByVkId(vkId int64) (entities.User, error) {
	user := entities.User{}
	err := s.db.Conn.QueryRow(s.db.Ctx, "SELECT id,vk_id,created_at FROM users WHERE vk_id=$1", vkId).
		Scan(&user.ID, &user.VkID, &user.CreatedAt)

	return user, err
}

func (s *UsersStorage) Create(user *entities.User) error {
	return s.db.Conn.
		QueryRow(s.db.Ctx, "INSERT INTO users(vk_id, created_at) VALUES ($1, $2) RETURNING id", user.VkID, time.Now()).
		Scan(&user.ID)
}
