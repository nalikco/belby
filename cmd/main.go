package main

import (
	"belby/internal/cli"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := cli.Handle(os.Args); err != nil {
		log.Fatal(err)
	}
}
