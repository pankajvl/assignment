package main

import (
	"fmt"
	"net/http"
)

func hellohandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
func main() {
	http.HandleFunc("/hello", hellohandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
