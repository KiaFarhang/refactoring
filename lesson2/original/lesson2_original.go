package original

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

/*
*
Imagine we have an application/piece of code that takes a client request, maps it to a command, and sends
that command for processing (e.g. to an event topic, batch process, whatever).

Now imagine this function was 200 lines long instead of ~10, with multiple dependency calls
and no unit tests - where do we even start refactoring?!?
*/
func (r *requestHandler) handleDeleteCreditCardRequest(request *deleteCreditCardRequest) error {
	/**
	This request-to-command mapping logic is a great candidate - can easily be pulled into a separate, pure function
	*/
	deleteCommand := &deleteCreditCardCommand{cardId: request.cardId}
	deleteCommand.cardholderName = fmt.Sprintf("%s, %s", request.cardholderLastName, request.cardholderFirstName)

	err := r.commandProcessor.DeleteCreditCard(deleteCommand)

	if err != nil {
		return err
	}

	return nil
}
