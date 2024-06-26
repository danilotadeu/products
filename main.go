package main

import (
	"log"

	serverInit "github.com/danilotadeu/products/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var (
	server serverInit.Server
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}
	server = serverInit.New()
}

func main() {
	server.Start()
}
