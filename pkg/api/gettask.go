package api

import (
	"net/http"
	"strconv"

	"github.com/th1nksnow/go-final-todo/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	var response taskResponse

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

	response.ID = task.ID
	response.Date = task.Date
	response.Title = task.Title
	response.Comment = task.Comment
	response.Repeat = task.Repeat

	writeJSON(w, response, http.StatusOK)
}
