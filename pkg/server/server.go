package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/th1nksnow/go-final-todo/pkg/api"
	"github.com/th1nksnow/go-final-todo/pkg/db"
)

type Config struct {
	Port   int
	webDir string
	dbFile string
}

func NewConfig() *Config {
	return &Config{
		Port:   7540,
		webDir: "./web",
		dbFile: "./scheduler.db",
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

func (c *Config) GetDbFile() string {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile != "" {
		c.dbFile = dbFile
	}

	return fmt.Sprintf(":%s", c.dbFile)
}

func NewServer() error {
	config := NewConfig()

	dbFile := config.GetDbFile()

	log.Printf("Initializing database %s", dbFile)
	if err := db.Init(config.dbFile); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	port := config.GetPort()

	api.Init()

	fileServer := http.FileServer(http.Dir(config.webDir))

	http.Handle("/", fileServer)

	log.Printf("Starting server on port %s", port)
	log.Printf("Serving files from: %s", config.webDir)

	return http.ListenAndServe(port, nil)
}
