package main

import (
	"io"
	"log"
	"net/http"
)

func getProducts(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hi!")
}

func main() {
	http.HandleFunc("/", getProducts)
	err := http.ListenAndServe(":4000", nil)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}
}
