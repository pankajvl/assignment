package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func idGenerator() func() int {
	id := 0
	return func() int {
		id++
		return id
	}
}

var (
	tasks     []Task
	getNextID = idGenerator()
)

func listTasksHandler(w http.ResponseWriter, r *http.Request) {
	var pending []Task
	for _, t := range tasks {
		if !t.Completed {
			pending = append(pending, t)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pending)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Description string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Description == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newTask := Task{
		ID:          getNextID(),
		Description: input.Description,
		Completed:   false,
	}
	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Completed = true
			fmt.Fprintf(w, "Task %d marked as completed\n", id)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listTasksHandler(w, r)
		case http.MethodPost:
			addTaskHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			completeTaskHandler(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
