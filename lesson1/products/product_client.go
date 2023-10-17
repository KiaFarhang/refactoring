package products

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
*
Define interfaces where they're used/needed, not where they're implemented.
https://github.com/golang/go/wiki/CodeReviewComments#interfaces
*/
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type productsClient struct {
	httpClient httpClient
}

/*
*
Expose a constructor that allows passing in any dependencies. We use this in unit tests to
provide a mocked httpClient, but it's useful outside of unit tests too - e.g. if we had one
HTTP client we were using throughout the app and wanted our products client to use it too.
*/
func NewCustomClient(httpClient httpClient) *productsClient {
	return &productsClient{httpClient: httpClient}
}

/*
*
Default constructor just provides sensible defaults for dependencies.
*/
func NewClient() *productsClient {
	return &productsClient{httpClient: &http.Client{Timeout: time.Second * 5}}
}

const productsEndpointUrl string = "https://dummyjson.com/products/%s"

func (p *productsClient) GetProduct(ctx context.Context, productId string) (*Product, error) {
	/**
	Note: in a real app all this "make a GET request, check the status code" HTTP logic would
	probably be in some separate, generic HTTP package that this client would just use - but you
	get the idea.
	*/
	url := fmt.Sprintf(productsEndpointUrl, productId)

	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	/**
	This error-check is not unit tested - time is money and I'd rather focus on
	logic that is more likely to go wrong/be coded incorrectly.
	*/
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
