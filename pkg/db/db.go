package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"` // формат YYYYMMDD
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

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

	DB, err = sql.Open("sqlite", dbFile)
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
	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}
	log.Println("Database schema created successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return fmt.Errorf("failed to close database")
}
