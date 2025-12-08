package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	password    string
	authEnabled bool
}

type taskResponse struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

var cfgApi *apiConfig

func taskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode JSON: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func initApiConfig() {
	// Дефолтные значения
	cfgApi = &apiConfig{
		password:    "",
		authEnabled: false,
	}

	authState := cfgApi.GetAuthPass()
	log.Printf("Initializing api config: %s", authState)
}

func (c *apiConfig) GetAuthPass() string {
	authPassword := os.Getenv("TODO_PASSWORD")
	if authPassword != "" {
		c.password = authPassword
		c.authEnabled = true
	}

	return fmt.Sprintf("authentication enabled: %t", c.authEnabled)
}

func Init() {
	initApiConfig()

	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))
}
