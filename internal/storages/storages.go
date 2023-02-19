package storages

import (
	"belby/internal/entities"
	"belby/pkg/database"
)

type Users interface {
	GetByVkId(vkId int64) (entities.User, error)
	Create(user *entities.User) error
}

type Messages interface {
	Create(message *entities.Message) error
}

type Products interface {
	Create(product *entities.Product) error
}

type Storages struct {
	Users
	Messages
	Products
}

func NewStorages(db *database.Database) *Storages {
	return &Storages{
		Users:    NewUsersStorage(db),
		Messages: NewMessagesStorage(db),
		Products: NewProductsStorage(db),
	}
}
