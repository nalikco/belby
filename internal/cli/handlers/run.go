package handlers

import (
	"belby/pkg/vk"
	"os"
)

func Run() error {
	vkApi := vk.NewVk(
		os.Getenv("VK_GROUP_ID"),
		os.Getenv("VK_TOKEN"),
	)

	if err := vkApi.Polling(); err != nil {
		return err
	}

	return nil
}
