package main

import (
	"fmt"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Starting server at port " + port)
	http.ListenAndServe(":"+port, nil)
}
