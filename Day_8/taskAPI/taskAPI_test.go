package main

// import (
//	"encoding/json"
//	"fmt"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
// )
// func TestViewTasks(t *testing.T) {
//	tests := []struct {
//		name           string
//		setupTasks     []Task
//		expectedStatus int
//		expectedCount  int
//	}{
//		{
//			name:           "Empty task list",
//			setupTasks:     []Task{},
//			expectedStatus: http.StatusOK,
//			expectedCount:  0,
//		},
//		{
//			name: "Single task",
//			setupTasks: []Task{
//				{ID: 1, Description: "Test task", Completed: false},
//			},
//			expectedStatus: http.StatusOK,
//			expectedCount:  1,
//		},
//		{
//			name: "Multiple tasks with mixed completion status",
//			setupTasks: []Task{
//				{ID: 1, Description: "Complete task", Completed: true},
//				{ID: 2, Description: "Incomplete task", Completed: false},
//				{ID: 3, Description: "Another task", Completed: true},
//			},
//			expectedStatus: http.StatusOK,
//			expectedCount:  3,
//		},
//		{
//			name: "Large number of tasks",
//			setupTasks: func() []Task {
//				tasks := make([]Task, 100)
//				for i := 0; i < 100; i++ {
//					tasks[i] = Task{
//						ID:          i + 1,
//						Description: fmt.Sprintf("Task %d", i+1),
//						Completed:   i%2 == 0, // Alternate between completed and incomplete
//					}
//				}
//				return tasks
//			}(),
//			expectedStatus: http.StatusOK,
//			expectedCount:  100,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			// Create a new server instance for each test
//			server := NewServer()
//
//			// Setup tasks
//			server.tasksMu.Lock()
//			for _, task := range tt.setupTasks {
//				server.tasksMap[task.ID] = &Task{
//					ID:          task.ID,
//					Description: task.Description,
//					Completed:   task.Completed,
//				}
//				// Update nextTaskID to ensure it's higher than existing IDs
//				if task.ID >= server.nextTaskID {
//					server.nextTaskID = task.ID
//				}
//			}
//			server.tasksMu.Unlock()
//
//			// Create request
//			req := httptest.NewRequest("GET", "/task", nil)
//			w := httptest.NewRecorder()
//
//			// Call the handler
//			server.viewTasks(w, req)
//
//			// Check status code
//			if w.Code != tt.expectedStatus {
//				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
//			}
//
//			// Check content type
//			expectedContentType := "application/json"
//			if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
//				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
//			}
//
//			// Parse response body
//			var responseTasks []Task
//			err := json.Unmarshal(w.Body.Bytes(), &responseTasks)
//			if err != nil {
//				t.Fatalf("Failed to unmarshal response: %v", err)
//			}
//
//			// Check number of tasks
//			if len(responseTasks) != tt.expectedCount {
//				t.Errorf("Expected %d tasks, got %d", tt.expectedCount, len(responseTasks))
//			}
//
//			// Verify all expected tasks are present (for non-empty cases)
//			if tt.expectedCount > 0 {
//				taskMap := make(map[int]Task)
//				for _, task := range responseTasks {
//					taskMap[task.ID] = task
//				}
//
//				for _, expectedTask := range tt.setupTasks {
//					actualTask, exists := taskMap[expectedTask.ID]
//					if !exists {
//						t.Errorf("Expected task with ID %d not found in response", expectedTask.ID)
//						continue
//					}
//
//					if actualTask.Description != expectedTask.Description {
//						t.Errorf("Task ID %d: expected description %q, got %q",
//							expectedTask.ID, expectedTask.Description, actualTask.Description)
//					}
//
//					if actualTask.Completed != expectedTask.Completed {
//						t.Errorf("Task ID %d: expected completed %t, got %t",
//							expectedTask.ID, expectedTask.Completed, actualTask.Completed)
//					}
//				}
//			}
//		})
//	}
// }
//
// func TestViewTasksConcurrency(t *testing.T) {
//	server := NewServer()
//
//	// Add some initial tasks
//	server.tasksMu.Lock()
//	for i := 1; i <= 10; i++ {
//		server.tasksMap[i] = &Task{
//			ID:          i,
//			Description: fmt.Sprintf("Task %d", i),
//			Completed:   false,
//		}
//		server.nextTaskID = i
//	}
//	server.tasksMu.Unlock()
//
//	// Test concurrent access to viewTasks
//	const numGoroutines = 50
//	done := make(chan bool, numGoroutines)
//	errors := make(chan error, numGoroutines)
//
//	for i := 0; i < numGoroutines; i++ {
//		go func() {
//			req := httptest.NewRequest("GET", "/task", nil)
//			w := httptest.NewRecorder()
//
//			server.viewTasks(w, req)
//
//			// Check that we got a valid response
//			if w.Code != http.StatusOK {
//				errors <- fmt.Errorf("Expected status 200, got %d", w.Code)
//				return
//			}
//
//			var tasks []Task
//			if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
//				errors <- fmt.Errorf("Failed to unmarshal response: %v", err)
//				return
//			}
//
//			if len(tasks) != 10 {
//				errors <- fmt.Errorf("Expected 10 tasks, got %d", len(tasks))
//				return
//			}
//
//			done <- true
//		}()
//	}
//
//	// Wait for all goroutines to complete
//	for i := 0; i < numGoroutines; i++ {
//		select {
//		case <-done:
//			// Success
//		case err := <-errors:
//			t.Errorf("Concurrency test failed: %v", err)
//		}
//	}
// }
//
// func TestViewTasksResponseFormat(t *testing.T) {
//	server := NewServer()
//
//	// Add a task with specific values to test JSON format
//	server.tasksMu.Lock()
//	server.tasksMap[1] = &Task{
//		ID:          1,
//		Description: "Test JSON format",
//		Completed:   true,
//	}
//	server.tasksMu.Unlock()
//
//	req := httptest.NewRequest("GET", "/task", nil)
//	w := httptest.NewRecorder()
//
//	server.viewTasks(w, req)
//
//	// Verify response is valid JSON array
//	var tasks []Task
//	err := json.Unmarshal(w.Body.Bytes(), &tasks)
//	if err != nil {
//		t.Fatalf("Response is not valid JSON: %v", err)
//	}
//
//	// Verify the structure matches expected format
//	if len(tasks) != 1 {
//		t.Fatalf("Expected 1 task, got %d", len(tasks))
//	}
//
//	task := tasks[0]
//	if task.ID != 1 {
//		t.Errorf("Expected ID 1, got %d", task.ID)
//	}
//	if task.Description != "Test JSON format" {
//		t.Errorf("Expected description 'Test JSON format', got %q", task.Description)
//	}
//	if task.Completed != true {
//		t.Errorf("Expected completed true, got %t", task.Completed)
//	}
//
//	// Verify JSON structure by checking raw bytes contain expected fields
//	bodyStr := w.Body.String()
//	expectedFields := []string{`"id":1`, `"description":"Test JSON format"`, `"completed":true`}
//	for _, field := range expectedFields {
//		if !contains(bodyStr, field) {
//			t.Errorf("Response JSON missing expected field: %s", field)
//		}
//	}
// }
//
// // Helper function to check if a string contains a substring
// func contains(s, substr string) bool {
//	return strings.Contains(s, substr)
// }
//
// // You'll need to add this import at the top:
// // "strings"
