package db

import "database/sql"

type Task struct {
	ID      int    `json:"id"`
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
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to insert task: %v", err)
	// }

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	// }

	// return id, nil
	if err == nil {
		id, err = result.LastInsertId()
	}
	return id, err
}
