package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KiaFarhang/refactoring/lesson1/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ProductClient interface {
	GetProduct(ctx context.Context, productId string) (*products.Product, error)
}

type server struct {
	productClient ProductClient
}

func NewServer(productClient ProductClient) *server {
	return &server{productClient: productClient}
}

/*
*
This is not production-ready code! You'd probably want to return custom error messages, for example.
The main point is that we're handling ONLY the web-layer logic here, and delegating the actual
"get the product" code to the client layer.
*/
func (s *server) GetProduct(w http.ResponseWriter, r *http.Request) {
	productIdParam := chi.URLParam(r, "productId")

	if len(productIdParam) == 0 {
		http.Error(w, http.StatusText(400), 400)
	}

	/**
	This Atoi and error check call is the ONLY piece of code that actually relates to our fictional
	piece of work - everything else was just refactoring the code to make it easy to add (and test).
	*/
	_, err := strconv.Atoi(productIdParam)

	if err != nil {
		http.Error(w, http.StatusText(400), 400)
	}

	product, err := s.productClient.GetProduct(r.Context(), productIdParam)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}

	if err := render.Render(w, r, product); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

}
