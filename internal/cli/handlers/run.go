package handlers

import (
	"belby/internal/handler"
	"belby/internal/storages"
	"belby/pkg/database"
	"belby/pkg/vk"
	"fmt"
	"os"
	"reflect"
)

func Run() error {
	vkApi := vk.NewVk(
		os.Getenv("VK_GROUP_ID"),
		os.Getenv("VK_TOKEN"),
	)

	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_DB_NAME"),
		SslMode:  os.Getenv("DB_SSL_MODE"),
	}
	db, err := database.NewDatabase(&dbConfig)
	if err != nil {
		return err
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	storage := storages.NewStorages(db)

	if err := vkApi.Polling(func(updates []interface{}, vkApi *vk.Vk) {
		for _, update := range updates {
			if reflect.TypeOf(update) == reflect.TypeOf(vk.Message{}) {
				if err := handler.HandleVk(storage, update.(vk.Message), vkApi); err != nil {
					fmt.Println(err)
				}
			}
		}
	}); err != nil {
		return err
	}

	return nil
}
