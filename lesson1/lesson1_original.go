package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/products/{productID}", func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("https://dummyjson.com/products/%s", chi.URLParam(r, "productID")), nil)
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

		var product Product

		err = json.Unmarshal(bodyAsBytes, &product)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		if err := render.Render(w, r, &product); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

	})
	err := http.ListenAndServe(":4000", r)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}
}

type Product struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (p *Product) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
