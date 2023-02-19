package handlers

import (
	"belby/internal/handler"
	"belby/pkg/vk"
	"os"
)

func Run() error {
	vkApi := vk.NewVk(
		os.Getenv("VK_GROUP_ID"),
		os.Getenv("VK_TOKEN"),
	)

	if err := vkApi.Polling(func(updates []vk.Update, vkApi *vk.Vk) {
		for _, update := range updates {
			if update.Type == "message_new" {
				err := handler.HandleVk(update.Object.(*vk.Message), vkApi)
				if err != nil {
					continue
				}
			}
		}
	}); err != nil {
		return err
	}

	return nil
}
