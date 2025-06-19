package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func hellohandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Im proper url")
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
	}

}

var m = make(map[int]string)
var i int = 0

func addTask(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	msg, _ := io.ReadAll(r.Body)
	m[i] = string(msg)
	i++

}

func getByID(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	url := r.URL.Path
	splited_path := strings.Split(url, "/")
	ind := splited_path[len(splited_path)-1]
	index, err := strconv.Atoi(ind)
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	w.Write([]byte(m[index]))
}

func main() {
	http.HandleFunc("/", hellohandler)

	http.HandleFunc("/task", addTask)
	http.HandleFunc("/task/{id}", getByID)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Not able to start server")
	}
}
