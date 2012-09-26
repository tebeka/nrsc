package main

import (
	"fmt"
	"net/http"
	"os"

	"./nrsc"
)

func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World\n")
}

func main() {
	nrsc.Handle("/static/")
	http.HandleFunc("/", indexHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
