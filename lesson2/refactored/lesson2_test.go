package refactored

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDeleteCreditCardCommand(t *testing.T) {
	t.Run("sets cardholderName to '{last name}, {first name}'", func(t *testing.T) {
		type test struct {
			request      *deleteCreditCardRequest
			expectedName string
		}

		tests := []test{
			{request: &deleteCreditCardRequest{cardholderFirstName: "Charli", cardholderLastName: "XCX"}, expectedName: "XCX, Charli"},
			{request: &deleteCreditCardRequest{cardholderFirstName: "Jon", cardholderLastName: "Snow"}, expectedName: "Snow, Jon"},
			/**
			Should we even be calling this function with requests that have empty cardholder names?
			Probably not - oftentimes unit testing small pieces of code like this exposes edge cases you
			need to discuss as a team.
			*/
			{request: &deleteCreditCardRequest{cardholderFirstName: "", cardholderLastName: ""}, expectedName: ", "},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("Given %s %s, should return %s", test.request.cardholderFirstName, test.request.cardholderLastName, test.expectedName), func(t *testing.T) {
				result := buildDeleteCreditCardCommand(test.request)
				assert.Equal(t, test.expectedName, result.cardholderName)
			})
		}
	})
}
