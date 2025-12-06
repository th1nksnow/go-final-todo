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

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`

	row := db.QueryRow(query, sql.Named("id", id))

	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`

	result, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))

	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if count == 0 {
		return errors.New("incorrect id for updating task")
	}

	return nil
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

// DeleteTask удаляет задачу по её ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`

	result, err := db.Exec(query, sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	// Проверяем, была ли удалена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if count == 0 {
		return errors.New("incorrect id for deleting task")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(id string, date string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`

	result, err := db.Exec(query,
		sql.Named("date", date),
		sql.Named("id", id))

	if err != nil {
		return fmt.Errorf("failed to update task date: %v", err)
	}

	// Проверяем, была ли обновлена хотя бы одна запись
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if count == 0 {
		return errors.New("incorrect id for updating task date")
	}

	return nil
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
