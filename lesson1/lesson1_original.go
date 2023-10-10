package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequestWithContext(r.Context(), "GET", "https://dummyjson.com/products", nil)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		httpClient := &http.Client{Timeout: time.Second * 5}
		response, err := httpClient.Do(request)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		defer response.Body.Close()

		statusCode := response.StatusCode

		if statusCode >= 400 {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		bodyAsBytes, err := io.ReadAll(response.Body)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.Write(bodyAsBytes)

	})
	err := http.ListenAndServe(":4000", r)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}
}
