package products

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

const productsEndpointUrl string = "https://dummyjson.com/products/%s"

func (p *productsClient) GetProduct(ctx context.Context, productId string) (*Product, error) {
	url := fmt.Sprintf(productsEndpointUrl, productId)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	response, err := p.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	statusCode := response.StatusCode

	if statusCode >= 400 {
		return nil, fmt.Errorf("received %d status code calling products API", statusCode)
	}

	bodyAsBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var product Product

	err = json.Unmarshal(bodyAsBytes, &product)

	if err != nil {
		return nil, err
	}

	return &product, nil
}
