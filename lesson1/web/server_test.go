package web

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KiaFarhang/refactoring/lesson1/products"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockProductClient struct {
	product *products.Product
	err     error
}

func (m *mockProductClient) GetProduct(ctx context.Context, productId string) (*products.Product, error) {
	return m.product, m.err
}

func TestGetProduct(t *testing.T) {
	t.Run("Returns a 400 if there's no 'productId' URL param in the request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/products/", nil)

		/**
		Set up a chi route context on the request, so the code can read it
		https://github.com/go-chi/chi/issues/76#issuecomment-370145140
		Though note we don't actually SET the 'productId' URL param - because that's
		the scenario we're testing
		*/
		requestContext := chi.NewRouteContext()

		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, requestContext))

		server := NewServer(&mockProductClient{})

		server.GetProduct(w, r)

		response := w.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
	t.Run("Returns a 400 if the product ID URL param isn't a number", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/products/", nil)

		requestContext := chi.NewRouteContext()
		// This time we actually do set the route param
		requestContext.URLParams.Add("productId", "foobar")

		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, requestContext))

		server := NewServer(&mockProductClient{})

		server.GetProduct(w, r)

		response := w.Result()

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	t.Run("Returns a 500 if there's an error fetching the product", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/products/", nil)

		requestContext := chi.NewRouteContext()
		requestContext.URLParams.Add("productId", "2")

		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, requestContext))

		// Tell our mock product client to return an error
		server := NewServer(&mockProductClient{err: errors.New("Some error calling the service")})

		server.GetProduct(w, r)

		response := w.Result()

		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	})
}
