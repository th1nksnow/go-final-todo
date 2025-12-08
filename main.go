package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/th1nksnow/go-final-todo/pkg/db"
	"github.com/th1nksnow/go-final-todo/pkg/server"
)

func setupSignalHandlers() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down server...")

		// Закрываем соединение с БД
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}

		log.Println("Server stopped gracefully")
		os.Exit(0)
	}()
}

func main() {
	log.Println("Starting TODO app...")

	// Ловим сигналы для graceful shutdown
	setupSignalHandlers()

	if err := server.NewServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
