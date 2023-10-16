package refactored

import "fmt"

type deleteCreditCardRequest struct {
	cardId              int
	cardholderFirstName string
	cardholderLastName  string
}

type deleteCreditCardCommand struct {
	cardId         int
	cardholderName string
}

type commandProcessor interface {
	DeleteCreditCard(command *deleteCreditCardCommand) error
}

type requestHandler struct {
	commandProcessor commandProcessor
}

func (r *requestHandler) handleDeleteCreditCardRequest(request *deleteCreditCardRequest) error {

	/**
	As a reader, the handleDeleteCreditCardRequest function is now easier to follow.
	All I see here is "oh, we pass in the request to build the command."

	If I WANT more detail I can jump into the mapping function - but I can stay at this layer if I choose.
	*/
	deleteCommand := buildDeleteCreditCardCommand(request)

	err := r.commandProcessor.DeleteCreditCard(deleteCommand)

	if err != nil {
		return err
	}

	return nil
}

/*
*
This is the only thing we changed in the refactored version - pull the mapping logic out and test it separately.
*/
func buildDeleteCreditCardCommand(request *deleteCreditCardRequest) *deleteCreditCardCommand {
	return &deleteCreditCardCommand{
		cardId:         request.cardId,
		cardholderName: fmt.Sprintf("%s, %s", request.cardholderLastName, request.cardholderFirstName),
	}
}
