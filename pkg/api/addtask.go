package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/th1nksnow/go-final-todo/pkg/db"
)

type addTaskRequest struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type addTaskResponse struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var response addTaskResponse
	var taskReq addTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		response.Error = fmt.Sprintf("Failed to decode JSON: %v", err)
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	task, err := validateAndProcessTask(taskReq)
	if err != nil {
		response.Error = err.Error()
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	id, err := db.AddTask(task)
	if err != nil {
		response.Error = fmt.Sprintf("Failed to add Task: %v", err)
		writeJSON(w, response, http.StatusInternalServerError)
		return
	}

	response.ID = int(id)
	writeJSON(w, response, http.StatusOK)
}

func validateAndProcessTask(taskReq addTaskRequest) (*db.Task, error) {
	if taskReq.Title == "" {
		return nil, errors.New("Не указан заголовок задачи")
	}

	now := time.Now()
	task := &db.Task{
		Title:   taskReq.Title,
		Comment: taskReq.Comment,
		Repeat:  taskReq.Repeat,
	}

	if err := processTaskDate(task, taskReq.Date, now); err != nil {
		return nil, err
	}

	return task, nil
}

func processTaskDate(task *db.Task, dateStr string, now time.Time) error {
	if dateStr == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return errors.New("Неверный формат даты. Ожидаемый формат: ГГГГММДД")
	}

	if task.Repeat != "" {
		// Проверяем корректность правила повторения и получаем nextDate
		nextDate, err := NextDate(now, dateStr, task.Repeat)
		if err != nil {
			return fmt.Errorf("Неверный формат правила повторения: %v", err)
		}

		// Если исходная дата меньше текущей, используем nextDate
		if afterNow(date, now) {
			task.Date = dateStr
		} else {
			task.Date = nextDate
		}
	} else {
		// Если правила повторения нет и дата меньше текущей, используем now
		if afterNow(date, now) {
			task.Date = dateStr
		} else {
			task.Date = now.Format(DateFormat)
		}
	}

	return nil
}
