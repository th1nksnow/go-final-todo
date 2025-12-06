package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/th1nksnow/go-final-todo/pkg/db"
)

// Обработчик POST запроса для отметки задачи выполненной
func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	var response taskResponse

	if r.Method != http.MethodPost {
		response.Error = "method not allowed"
		writeJSON(w, response, http.StatusMethodNotAllowed)
		return
	}

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

	task, err := db.GetTask(id)
	if err != nil {
		response.Error = err.Error()
		writeJSON(w, response, http.StatusInternalServerError)
		return
	}

	now := time.Now()

	// Проверяем, является ли задача одноразовой
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			response.Error = fmt.Sprintf("error deleting task: %v", err)
			writeJSON(w, response, http.StatusInternalServerError)
			return
		}
	} else {
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			response.Error = fmt.Sprintf("error processing next date: %v", err)
			writeJSON(w, response, http.StatusInternalServerError)
			return
		}

		err = db.UpdateDate(id, nextDate)
		if err != nil {
			response.Error = fmt.Sprintf("error updating task date: %v", err)
			writeJSON(w, response, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{}, http.StatusOK)
}
