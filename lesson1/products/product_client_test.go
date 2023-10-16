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

/*
*
Create a mock implementation of the interface our client under test depends on
Sometimes mocks need more complicated logic, but often you can get away with
"just return whatever I tell you when the method is called" like this.
*/
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
		/**
		Notice how we create a new mock HTTP client and client under test in every unit test
		This is intentional; I usually prefer the extra lines of boilerplate because each test is
		self-contained and includes everything you need to see at a glance.

		Vs. having one mock that each test shares, overwrites, etc. - can get messy fast.
		*/
		productClient := NewCustomClient(&mockHttpClient{mockError: mockError})
		product, err := productClient.GetProduct(context.Background(), productId)
		assert.Nil(t, product)
		assert.Error(t, err)
	})
	t.Run("Returns an error if the status code is >= 400", func(t *testing.T) {
		statusCodes := []int{400, 401, 403, 500, 503}
		/**
		Table-driven tests let us run many tests with changing parameters without actually
		having to WRITE each test separately.
		https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
		*/
		for _, statusCode := range statusCodes {
			t.Run(fmt.Sprintf("Status code %d", statusCode), func(t *testing.T) {
				mockResponse := buildMockHttpResponse(statusCode, "")
				productClient := NewCustomClient(&mockHttpClient{mockResponse: mockResponse})
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
		productClient := NewCustomClient(&mockHttpClient{mockResponse: mockResponse})
		product, err := productClient.GetProduct(context.Background(), productId)
		assert.Nil(t, product)
		assert.Error(t, err)
	})

	t.Run("Returns a Product if the response body can be parsed", func(t *testing.T) {
		json := `{"id": 123, "title": "iPhone"}`
		mockResponse := buildMockHttpResponse(200, json)
		productClient := NewCustomClient(&mockHttpClient{mockResponse: mockResponse})
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
