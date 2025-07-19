package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/juliotorresmoreno/odoo-executor/handler"
)

func main() {
	_ = godotenv.Load(".env")

	handler := handler.ConfigureHandler()

	httpServer := http.Server{
		Addr:    ":4080",
		Handler: handler,
	}

	log.Println(httpServer.ListenAndServe())
}
