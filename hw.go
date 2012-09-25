package main

import (
	"fmt"
	"net/http"

	"./resources"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	resources.Handle("/static/")
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
