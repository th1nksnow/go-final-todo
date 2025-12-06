package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/th1nksnow/go-final-todo/pkg/db"
)

type updateTaskRequest struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Обработчик PUT запроса для обновления задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var response taskResponse
	var taskReq updateTaskRequest

	id := r.FormValue("id")
	if id == "" {
		response.Error = "id not specified"
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		response.Error = "invalid id format"
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		response.Error = fmt.Sprintf("failed to decode JSON: %v", err)
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	task, err := validateAndProcessUpdateTask(taskReq)
	if err != nil {
		response.Error = err.Error()
		writeJSON(w, response, http.StatusBadRequest)
		return
	}

	// Обновляем задачу в БД
	err = db.UpdateTask(task)
	if err != nil {
		response.Error = err.Error()
		writeJSON(w, response, http.StatusNotFound)
		return
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}

// Функция для валидации и обработки данных при обновлении задачи
func validateAndProcessUpdateTask(taskReq updateTaskRequest) (*db.Task, error) {
	if taskReq.Title == "" {
		return nil, errors.New("the title is not specified")
	}
	now := time.Now()
	task := &db.Task{
		ID:      taskReq.ID,
		Title:   taskReq.Title,
		Comment: taskReq.Comment,
		Repeat:  taskReq.Repeat,
	}

	if err := processTaskDate(task, taskReq.Date, now); err != nil {
		return nil, err
	}

	return task, nil
}
