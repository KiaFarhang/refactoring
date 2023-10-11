package products

import "net/http"

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
