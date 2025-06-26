package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func hellohandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello, World!")
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	}

}

type Record struct {
	task      string
	completed bool
}

var m = make(map[int]*Record)
var i int = 0

func addTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	msg, _ := io.ReadAll(r.Body)
	m[i] = &Record{string(msg), false}
	i++

}

func getByID(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	w.Write([]byte(strconv.Itoa(index)))
	w.WriteHeader(200)
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	for _, task := range m {
		fmt.Fprintf(w, "%v\n", task)
	}
}

func completeTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	t, ok := m[index]
	if !ok {
		fmt.Fprintf(w, "%v is not found\n", index)
		w.WriteHeader(404)
		return
	}

	t.completed = true
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	delete(m, index)
}

func main() {
	http.HandleFunc("/", hellohandler)

	http.HandleFunc("POST /task", addTask)
	http.HandleFunc("GET /task/{id}", getByID)
	http.HandleFunc("GET /task", viewTask)
	http.HandleFunc("PUT /task/{id}", completeTask)
	http.HandleFunc("DELETE /task/{id}", deleteTask)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Not able to start server")
	}
}
