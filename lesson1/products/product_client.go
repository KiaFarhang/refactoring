package products

import (
	"context"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ProductsClient interface {
	GetProduct(ctx context.Context, productId string) (*Product, error)
}

type productsClient struct {
	httpClient httpClient
}

func NewProductsClient(httpClient httpClient) ProductsClient {
	return &productsClient{httpClient: httpClient}
}

func (p *productsClient) GetProduct(ctx context.Context, productId string) (*Product, error) {
	return &Product{}, nil
}
