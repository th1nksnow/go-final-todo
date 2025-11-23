package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Port   int
	WebDir string
}

func NewConfig() *Config {
	return &Config{
		Port:   7540,
		WebDir: "./web",
	}
}

func (c *Config) GetPort() string {
	portStr := os.Getenv("TODO_PORT")
	if portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 && port < 65536 {
			c.Port = port
		} else {
			log.Printf("Invalid TODO_PORT value: %s, using default port: %d", portStr, c.Port)
		}
	}
	return fmt.Sprintf(":%d", c.Port)
}

func NewServer() error {
	config := NewConfig()
	port := config.GetPort()

	fileServer := http.FileServer(http.Dir(config.WebDir))

	http.Handle("/", fileServer)

	log.Printf("Starting server on port %s", port)
	log.Printf("Serving files from: %s", config.WebDir)

	return http.ListenAndServe(port, nil)
}
