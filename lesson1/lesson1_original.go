package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hi!"))
	})
	err := http.ListenAndServe(":4000", r)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}
}
