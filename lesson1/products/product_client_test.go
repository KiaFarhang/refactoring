package products

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockHttpClient struct {
	mockResponse http.Response
	mockError    error
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return &m.mockResponse, m.mockError
}

const productId string = "123"

func TestGetProduct(t *testing.T) {
	t.Run("Returns an error if the HTTP call fails", func(t *testing.T) {
		mockError := errors.New("Something went wrong!")
		productClient := &productsClient{httpClient: &mockHttpClient{mockError: mockError}}
		product, err := productClient.GetProduct(context.Background(), productId)
		assert.Nil(t, product)
		assert.Error(t, err)
	})
	t.Run("Returns an error if the status code is >= 400", func(t *testing.T) {
		statusCodes := []int{400, 401, 403, 500, 503}
		for _, statusCode := range statusCodes {
			t.Run(fmt.Sprintf("Status code %d", statusCode), func(t *testing.T) {
				mockResponse := buildMockHttpResponse(statusCode, "")
				productClient := &productsClient{httpClient: &mockHttpClient{mockResponse: mockResponse}}
				product, err := productClient.GetProduct(context.Background(), productId)
				assert.Nil(t, product)
				assert.Error(t, err)
				expectedErrorMessage := fmt.Sprintf("received %d status code calling products API", statusCode)
				assert.Equal(t, expectedErrorMessage, err.Error())
			})
		}
	})
	t.Run("Returns an error if the response body can't be parsed", func(t *testing.T) {
		mockResponse := buildMockHttpResponse(200, "Not valid JSON!")
		productClient := &productsClient{httpClient: &mockHttpClient{mockResponse: mockResponse}}
		product, err := productClient.GetProduct(context.Background(), productId)
		assert.Nil(t, product)
		assert.Error(t, err)
	})

	t.Run("Returns a Product if the response body can be parsed", func(t *testing.T) {
		json := `{"id": 123, "title": "iPhone"}`
		mockResponse := buildMockHttpResponse(200, json)
		productClient := &productsClient{httpClient: &mockHttpClient{mockResponse: mockResponse}}
		product, err := productClient.GetProduct(context.Background(), productId)
		assert.Nil(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 123, product.Id)
		assert.Equal(t, "iPhone", product.Title)
	})
}

func buildMockHttpResponse(statusCode int, jsonBody string) http.Response {
	return http.Response{StatusCode: statusCode, Body: io.NopCloser(strings.NewReader(jsonBody))}
}
