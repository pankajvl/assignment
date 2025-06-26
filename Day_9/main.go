package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	Completed   bool
}
type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}
func hellohandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}

func (s *Server) addTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqBody struct {
		Description string
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || json.Unmarshal(bodyBytes, &reqBody) != nil || reqBody.Description == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest) //400
		return
	}
	res, err := s.db.Exec("INSERT INTO tasks (description,completed) VALUES(?,?)", reqBody.Description, false)
	if err != nil {
		http.Error(w, "DB insert failed", http.StatusInternalServerError) //500
		return
	}
	id, _ := res.LastInsertId()
	task := Task{ID: int(id), Description: reqBody.Description, Completed: false}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
func (s *Server) getById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var task Task
	err = s.db.QueryRow("SELECT id,description,completed FROM tasks WHERE id=?", id).Scan(&task.ID, &task.Description, &task.Completed)
	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(task)
}

func (s *Server) viewTasks(w http.ResponseWriter, _ *http.Request) {
	rows, err := s.db.Query("SELECT id,description,completed FROM tasks")
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Description, &task.Completed); err == nil {
			tasks = append(tasks, task)
		}

	}
	json.NewEncoder(w).Encode(tasks)
}
func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invlid ID ", http.StatusBadRequest)
		return
	}
	res, err := s.db.Exec("UPDATE tasks SET completed=true WHERE id=?", id)
	if err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	s.getById(w, r)
}
func (s *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invlid ID ", http.StatusBadRequest)
		return
	}
	res, err := s.db.Exec("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent) //204
}
func main() {
	dsn := "root:root123@tcp(localhost:3306)/test_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB connection failed:%v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}
	server := NewServer(db)
	http.HandleFunc("/", hellohandler)
	http.HandleFunc("POST /task", server.addTask)
	http.HandleFunc("GET /task/{id}", server.getById)
	http.HandleFunc("GET /task", server.viewTasks)
	http.HandleFunc("PUT /task/{id}", server.completeTask)
	http.HandleFunc("DELETE /task/{id}", server.deleteTask)

	srv := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Starting server on :8000")
	log.Fatal(srv.ListenAndServe())
}
