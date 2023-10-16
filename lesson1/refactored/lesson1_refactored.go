package main

import (
	"log"
	"net/http"

	"github.com/KiaFarhang/refactoring/lesson1/products"
	"github.com/KiaFarhang/refactoring/lesson1/web"
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

	// Instantiate our client without passing in a dependency, letting it construct a default one internally
	productClient := products.NewClient()
	server := web.NewServer(productClient)
	r.Get("/products/{productId}", server.GetProduct)
	err := http.ListenAndServe(":4000", r)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}

}
