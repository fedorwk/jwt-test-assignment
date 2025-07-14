package main

import (
	"log"
	"medods-auth/app/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := server.Start(); err != nil {
		log.Fatalln(err)
	}
}
