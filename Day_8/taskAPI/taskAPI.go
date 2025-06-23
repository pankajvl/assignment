package main

import (
	"encoding/json" // For JSON encoding/decoding
	"fmt"
	"io"  // For reading request body
	"log" // For server logs
	"net/http"
	"strconv" // For converting ID strings to integers
	"sync"    // For mutex to handle concurrent access to shared data
	"time"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type Server struct {
	tasksMap   map[int]*Task
	nextTaskID int
	tasksMu    sync.Mutex // Protects tasksMap and nextTaskID
}

func NewServer() *Server {
	return &Server{
		tasksMap:   make(map[int]*Task),
		nextTaskID: 0,
	}
}

func (s *Server) getNextID() int {
	s.tasksMu.Lock()
	defer s.tasksMu.Unlock()

	s.nextTaskID++

	return s.nextTaskID
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "Hello, Simplified Go Task API!")
}

func (s *Server) addTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var reqBody struct {
		Description string `json:"description"`
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(bodyBytes, &reqBody)

	if err != nil {
		http.Error(w, "Invalid JSON format or missing 'description'.", http.StatusBadRequest)
		return
	}

	if reqBody.Description == "" {
		http.Error(w, "Task description cannot be empty.", http.StatusBadRequest)
		return
	}

	newTask := Task{
		ID:          s.getNextID(),
		Description: reqBody.Description,
		Completed:   false,
	}

	s.tasksMu.Lock()
	s.tasksMap[newTask.ID] = &newTask
	s.tasksMu.Unlock()

	responseJSON, err := json.Marshal(newTask)
	if err != nil {
		http.Error(w, "Failed to prepare response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created

	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("addTask: Error writing response: %v", err)
	}
}

// getByID: Handles GET /task/{id} to retrieve a single task.
func (s *Server) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID format.", http.StatusBadRequest)
		return
	}

	s.tasksMu.Lock()
	task, ok := s.tasksMap[id]
	s.tasksMu.Unlock()

	if !ok {
		http.Error(w, fmt.Sprintf("Task with ID %d not found.", id), http.StatusNotFound)
		return
	}

	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to prepare response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("addTask: Error writing response: %v", err)
	}
}

func (s *Server) viewTasks(w http.ResponseWriter, _ *http.Request) {
	s.tasksMu.Lock()

	allTasks := []Task{}
	for _, task := range s.tasksMap {
		allTasks = append(allTasks, *task)
	}
	s.tasksMu.Unlock()

	responseJSON, err := json.Marshal(allTasks)
	if err != nil {
		http.Error(w, "Failed to prepare response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("addTask: Error writing response: %v", err)
	}
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID format.", http.StatusBadRequest)
		return
	}

	s.tasksMu.Lock()
	task, ok := s.tasksMap[id]

	if !ok {
		s.tasksMu.Unlock()
		http.Error(w, fmt.Sprintf("Task with ID %d not found.", id), http.StatusNotFound)

		return
	}

	task.Completed = true
	s.tasksMu.Unlock()

	responseJSON, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "Failed to prepare response.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK

	if _, err := w.Write(responseJSON); err != nil {
		log.Printf("addTask: Error writing response: %v", err)
	}
}

func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID format.", http.StatusBadRequest)
		return
	}

	s.tasksMu.Lock()
	_, ok := s.tasksMap[id]

	if !ok {
		s.tasksMu.Unlock()
		http.Error(w, fmt.Sprintf("Task with ID %d not found.", id), http.StatusNotFound)

		return
	}

	delete(s.tasksMap, id)
	s.tasksMu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	server := NewServer()

	id1 := server.getNextID()
	id2 := server.getNextID()

	server.tasksMu.Lock()
	server.tasksMap[id1] = &Task{ID: id1, Description: "Simplify code", Completed: false}
	server.tasksMap[id2] = &Task{ID: id2, Description: "Learn Go basics", Completed: false}
	server.tasksMu.Unlock()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("POST /task", server.addTask)
	http.HandleFunc("GET /task/{id}", server.getByID)
	http.HandleFunc("GET /task", server.viewTasks)
	http.HandleFunc("PUT /task/{id}", server.completeTask)
	http.HandleFunc("DELETE /task/{id}", server.deleteTask)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Simplified API Server starting on :8080...")
	log.Fatal(srv.ListenAndServe())
}
