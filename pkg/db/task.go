package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // формат YYYYMMDD
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`

	result, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = result.LastInsertId()
	}
	return id, err
}

func Tasks(search string, limit int) ([]*Task, error) {
	var query string
	var rows *sql.Rows
	var err error

	if limit < 1 || limit > 50 {
		return nil, errors.New("records limit should be in range from 1 to 50")
	}

	if search != "" {
		if parsedDate, err := time.Parse("02.01.2006", search); err == nil {
			dateStr := parsedDate.Format(DateFormat)
			query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE date = :date
	          ORDER BY date LIMIT :limit`
			rows, err = db.Query(query,
				sql.Named("date", dateStr),
				sql.Named("limit", limit))
		} else {
			query := `SELECT id, date, title, comment, repeat FROM scheduler
	          WHERE title LIKE :search OR comment LIKE :search
	          ORDER BY date LIMIT :limit`
			rows, err = db.Query(query,
				sql.Named("search", "%"+search+"%"),
				sql.Named("limit", limit))
		}
	} else {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit"
		rows, err = db.Query(query,
			sql.Named("limit", limit))
	}
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %v", err)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}
