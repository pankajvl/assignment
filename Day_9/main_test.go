package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setUpTestServer(t *testing.T) *Server {
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	_, err = db.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("Failed to clean up test DB: %v", err)
	}
	return NewServer(db)
}
func TestAddTask(t *testing.T) {
	s := setUpTestServer(t)
	payload := map[string]string{"description": "test"}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.addTask)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d but got %d", http.StatusCreated, status)
	}
	var task Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	if task.Description != "test" {
		t.Errorf("Expected task description %s but got %s", "test", task.Description)
	}
}

func TestGetTaskByID(t *testing.T) {
	s := setUpTestServer(t)

	res, err := s.db.Exec("INSERT INTO tasks (description, completed) VALUES (?, ?)", "get task test", false)
	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}
	id, _ := res.LastInsertId()

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%d", id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", id))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.getById)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}
	if task.ID != int(id) {
		t.Errorf("Expected ID %d, got %d", id, task.ID)
	}
}

func TestViewAllTasks(t *testing.T) {
	s := setUpTestServer(t)

	s.db.Exec("INSERT INTO tasks (description, completed) VALUES (?, ?)", "Task 1", false)
	s.db.Exec("INSERT INTO tasks (description, completed) VALUES (?, ?)", "Task 2", true)

	req := httptest.NewRequest(http.MethodGet, "/task", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.viewTasks)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var tasks []Task
	if err := json.Unmarshal(rr.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Failed to parse tasks: %v", err)
	}
	if len(tasks) < 2 {
		t.Errorf("Expected at least 2 tasks, got %d", len(tasks))
	}
}

func TestCompleteTask(t *testing.T) {
	s := setUpTestServer(t)

	res, _ := s.db.Exec("INSERT INTO tasks (description, completed) VALUES (?, ?)", "to complete", false)
	id, _ := res.LastInsertId()

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/task/%d", id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", id))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.completeTask)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var task Task
	if err := json.Unmarshal(rr.Body.Bytes(), &task); err != nil {
		t.Fatalf("Failed to parse task: %v", err)
	}
	if !task.Completed {
		t.Errorf("Expected task to be completed")
	}
}

func TestDeleteTask(t *testing.T) {
	s := setUpTestServer(t)

	res, _ := s.db.Exec("INSERT INTO tasks (description, completed) VALUES (?, ?)", "to delete", false)
	id, _ := res.LastInsertId()

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/task/%d", id), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", id))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.deleteTask)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rr.Code)
	}

	var count int
	row := s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ?", id)
	err := row.Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query deleted task: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected task to be deleted, but found it in DB")
	}
}

// ........
func TestAddTask_InvalidJSON(t *testing.T) {
	s := setUpTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.addTask)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid JSON, got %d", rr.Code)
	}
}
func TestAddTask_EmptyDescription(t *testing.T) {
	s := setUpTestServer(t)
	payload := map[string]string{"description": ""}
	payloadBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/task", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.addTask)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for empty description, got %d", rr.Code)
	}
}
func TestGetTaskByID_InvalidID(t *testing.T) {
	s := setUpTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/task/abc", nil)
	req.SetPathValue("id", "abc") // invalid ID
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(s.getById)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 for invalid ID, got %d", rr.Code)
	}
}

// ......
func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(hellohandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "Hello world!" {
		t.Errorf("Unexpected response body: %s", rr.Body.String())
	}
}
