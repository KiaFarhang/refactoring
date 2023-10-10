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

var httpClient = &http.Client{Timeout: time.Second * 5}

func main() {
	/**
	Set up a boilerplate HTTP router with chi: https://github.com/go-chi/chi
	Not super relevant for our purposes
	*/
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/products/{productID}", getProducts)
	err := http.ListenAndServe(":4000", r)

	if err != nil {
		log.Fatalf("Unable to start server, error was %v", err)
	}
}

/*
*
Fetches a given product from the Dummy JSON products API, taking the product ID from the URL of the request (e.g. someone hitting us at /products/5)
*/
func getProducts(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	url := fmt.Sprintf("https://dummyjson.com/products/%s", productID)
	request, err := http.NewRequestWithContext(r.Context(), "GET", url, nil)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

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
}

/*
*
Represents a product as returned from Dummy JSON's /products/{id} endpoint:
https://dummyjson.com/docs/products

We take just a couple fields from the upstream endpoint's response
*/
type Product struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

/*
*
A quirk of chi - this lets us easily render a Product as JSON when we want to send it back to our clients
*/
func (p *Product) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
