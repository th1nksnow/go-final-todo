package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// Init инициализирует бдазу данных
func Init(dbFile string) error {
	var install bool

	// Проверяем наличие файла с базой данных
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Если БД не существовала, создаем схему
	if install {
		log.Printf("Creating new database: %s", dbFile)
		if err := createSchema(); err != nil {
			return fmt.Errorf("failed to create schema: %v", err)
		}
	} else {
		log.Printf("Using existing database: %s", dbFile)
	}

	return nil
}

// createSchema создает таблицы и индексы
func createSchema() error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}
	log.Println("Database schema created successfully")
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return errors.New("failed to close database")
}
