package lesson3

import "testing"

func TestBatchJobRunner(t *testing.T) {
	/**
	This test doesn't actually...test anything. It calls the main() function, which swallows any errors,
	and doesn't assert anything about what happens afterwards.

	But hey - code coverage says we've now fully tested the main() function!

	This is why it's critical to review unit tests just as closely as source code - otherwise you'd
	have a false sense of security from a useless test.

	P.S.: I found this shell script really useful for quickly viewing Go test coverage:
	https://stackoverflow.com/a/27284510
	*/
	t.Run("Is a totally useless unit test", func(t *testing.T) {
		main()
	})
}
