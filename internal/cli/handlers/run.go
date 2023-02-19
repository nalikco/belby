package handlers

import (
	"belby/internal/handler"
	"belby/pkg/vk"
	"os"
	"reflect"
)

func Run() error {
	vkApi := vk.NewVk(
		os.Getenv("VK_GROUP_ID"),
		os.Getenv("VK_TOKEN"),
	)

	if err := vkApi.Polling(func(updates []interface{}, vkApi *vk.Vk) {
		for _, update := range updates {
			if reflect.TypeOf(update) == reflect.TypeOf(vk.Message{}) {
				_ = handler.HandleVk(update.(vk.Message), vkApi)
			}
		}
	}); err != nil {
		return err
	}

	return nil
}
