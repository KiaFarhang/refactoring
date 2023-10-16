package web

import (
	"context"
	"net/http"
	"strconv"

	"github.com/KiaFarhang/refactoring/lesson1/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

/*
*
Define interfaces where they're used/needed, not where they're implemented.
https://github.com/golang/go/wiki/CodeReviewComments#interfaces
*/
type ProductClient interface {
	GetProduct(ctx context.Context, productId string) (*products.Product, error)
}

type server struct {
	productClient ProductClient
}

func NewServer(productClient ProductClient) *server {
	return &server{productClient: productClient}
}

func (s *server) GetProduct(w http.ResponseWriter, r *http.Request) {
	productIdParam := chi.URLParam(r, "productId")

	if len(productIdParam) == 0 {
		http.Error(w, http.StatusText(400), 400)
	}

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
