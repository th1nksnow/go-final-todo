package api

import (
	"net/http"
	"strings"

	"github.com/th1nksnow/go-final-todo/pkg/db"
)

type TasksResp struct {
	Tasks *[]*db.Task `json:"tasks,omitempty"`
	Error string      `json:"error,omitempty"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	var response TasksResp

	if r.Method != http.MethodGet {
		response.Error = "method not allowed"
		writeJSON(w, response, http.StatusMethodNotAllowed)
		return
	}

	search := strings.TrimSpace(r.FormValue("search"))

	tasks, err := db.Tasks(search, 50)
	if err != nil {
		response.Error = "failed to get tasks: " + err.Error()
		writeJSON(w, response, http.StatusInternalServerError)
		return
	}

	// WARNING! If you return an error further, the response body will contain {"tasks":[], ...}
	response.Tasks = &tasks

	writeJSON(w, response, http.StatusOK)
}
