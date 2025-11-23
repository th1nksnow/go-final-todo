package main

import (
	"log"

	"github.com/th1nksnow/go-final-todo/pkg/server"
)

func main() {
	log.Println("Starting TODO app...")

	if err := server.NewServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
