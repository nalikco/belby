package main

import (
	"belby/pkg/vk"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	vkApi := vk.NewVk(
		os.Getenv("VK_GROUP_ID"),
		os.Getenv("VK_TOKEN"),
	)

	err = vkApi.Polling()
	if err != nil {
		log.Fatal(err)
	}
}
