# Lesson 1 - Separating Concerns and Dependency Injection

In this example, we have a small web application that provides details about products when users hit it with a GET on a product ID. The app just looks up products using [a dummy JSON API](https://dummyjson.com/docs/products).

You can run the original code and test it to see it does work:

```shell
go run lesson1_original.go
# In a separate terminal...
curl http://localhost:4000/products/2 | jq # the pipe to jq is optional
```

You should see a JSON representation of a dummy product.

Now, imagine you've been asked to add a new piece of functionality to this app. We want to return a 400 HTTP status code if the client calls us with a product ID that's not an integer - e.g. `http://localhost:4000/products/foobar`.

We _could_ add this code to the existing implementation - but it's a bit messy and there are no unit tests.

Instead, we're going to **use this opportunity to refactor as we add new functionality.** I find this is the best time to refactor/improve code - you're already in there, you don't need to sell anyone on a "make the code better but don't add any new functionality" ticket, etc.

Now, this is going to be team-specific: always ask before you take on a refactor like this. Maybe there isn't time right now, maybe we're deprecating this app next week so it doesn't matter.

But for our purposes, let's assume everyone's on board with you doing this refactoring as part of your work. And let's do it!

## Refactoring this code

Our refactoring focuses on a couple of areas:

First, we separate the pieces of functionality in this application - serving a web application and looking up products - into separate packages (`products` and `web`). 

In a dummy app like this that's overkill. But this is one of the easiest ways to start improving legacy code: Find a piece of functionality and split it out. This will help us make other improvements, like...

### Unit testing with dependency injection

#### Refactoring the API client

Let's take a look at the original API-calling code:

```go
func getProducts(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	url := fmt.Sprintf("https://dummyjson.com/products/%s", productID)
	request, err := http.NewRequestWithContext(r.Context(), "GET", url, nil)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	response, err := httpClient.Do(request)
    // Omitting the rest for brevity's sake
}
```

You'll see the code to call the dummy API uses a global HTTP client defined earlier in the file. It also relies directly on our web framework to grab the ID of the product the caller wants to look up - even though the web framework code is, ideally, completely divorced from this API client.

(For example, what if we wanted to trigger API calls from a batch job instead of an HTTP request? Or from a CLI command?)

This code as written is fairly hard to unit test. We'd need to find a way to simulate an HTTP GET so we can grab the product ID. We'd probably want to mock the dummy JSON API or set up some kind of proxy so our tests don't trigger real requests to it.

To make this easier, when refactoring we make a list of all the dependencies this API-calling code needs to do its job. Luckily, in this case that's just one thing: an HTTP client.

In [`products/product_client.go`](./products/product_client.go)We define the product client as a struct with a single field: an interface matching the _real_ HTTP client functionality we need:

```go
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type productsClient struct {
	httpClient httpClient
}
```

Then we offer two ways to create a `productsClient`:

```go
func NewCustomClient(httpClient httpClient) *productsClient {
	return &productsClient{httpClient: httpClient}
}

func NewClient() *productsClient {
	return &productsClient{httpClient: &http.Client{Timeout: time.Second * 5}}
}
```

The first function, `NewCustomClient`, is what allows us to inject a _fake_ HTTP client in [the unit tests for this client code.](./products/product_client_test.go) Now it's trivial to test all the different paths: What happens if the HTTP call fails? What happens if we get a status code we didn't expect? etc. We don't need to write brittle unit tests that actually result in calls to the underlying API.

The second function, `NewClient()`, is what the refactored code will actually use. It just creates the same HTTP client we did in the original version.

#### Refactoring the web server

We can follow a similar pattern to refactor and unit test the web server code to what lives in [`web/server.go`](./web/server.go).

First, identify all the dependencies our web server needs. In this case, that's just "a way to look up products."

```go
type ProductClient interface {
	GetProduct(ctx context.Context, productId string) (*products.Product, error)
}

type server struct {
	productClient ProductClient
}
```

Note that there's no mention of the dummy API (or even HTTP) in our `ProductClient` interface. That's because **the web server code shouldn't care how products are fetched** - it just knows it needs to ask for a product by ID, and get either that product or an error back.

Because we delegate to our product client dependency in the web server code:

```go
product, err := s.productClient.GetProduct(r.Context(), productIdParam)
```

We can similarly write [unit tests](./web/server_test.go) for the web server code that create a mock product client doing whatever we need it to do to exercise web server functionality.

### Pausing to PR what's there

At this point we've pretty drastically changed the organization of the application. We also added a bunch of unit tests to verify existing functionality.

This is when I'd generally recommend opening a pull request. It's another form of separating concerns: Someone can review _just_ your refactor + addition of test cases without _also_ having to think about the functionality you added. Your team could even release the refactored code to production and make sure there are no unexpected issues before you move ahead.

Keep your pull requests as small and focused as possible. It makes reviewing and troubleshooting code far simpler.

### Adding our new functionality (with tests!)

Let's assume our refactor PR merged. That's great, but we still need to do the thing we came for: return a 400 if the web server gets a non-integer product ID.

The great thing about refactoring legacy code is **it makes adding new functionality easy.** Because the code is separated into logical pieces, and each piece already has unit tests, our work is literally just this in the web server:

```go
	_, err := strconv.Atoi(productIdParam)

	if err != nil {
		http.Error(w, http.StatusText(400), 400)
	}
```

And a couple of new tests in [`web/server_test.go`](./web/server_test.go) to verify that functionality.

Our refactoring paid off - not just for us, but for every future developer who works on this application!

### A note on test-driven development

Throughout this walkthrough I typically referred first to source code, then the unit tests verifying it. But whenever possible I recommend you use test-driven development and do the opposite: Write unit tests describing what your code should do, _then_ write the code to actually do it.

I've found TDD helps in the following ways:

- It helps you define your dependencies and utility/helper code. By breaking your high-level problem into "it must do this," "under these circumstances it must do that," you'll oftentimes look at the resulting list and go "Ah - well, that piece of work can be its own function. That piece of work should be delegated to a dependency."
- On larger pieces of work, it can also help you divide and conquer. If you've written 10 tests across two different files/packages, great - one developer can get the tests in package A passing, and another can tackle package B simultaneously. This is another benefit of using interfaces for dependencies. If you define interfaces up front, you can implement those interfaces in parallel.
- TDD often helps you think of edge cases you wouldn't otherwise. I find that if I write tests after source code, I'm more likely to only test directly what's already there. With TDD, you can step back and think about all the edge cases and error scenarios and make sure your code handles them.
- TDD often leads to far simpler code. You've already broken your work into separate, unit-tested pieces. Now reviewing that code, updating it six months from now, etc. should be a breeze.

Now, there are cases where TDD isn't really appropriate:

- If you're building a proof-of-concept or other throwaway/temporary work
- If something is actively on fire and you just need to push code to resolve an incident (though definitely follow back up with tests for the fix!)
- If your tests aren't actually valuable/testing what you think they are (more on this in lesson 3)

I try not to be dogmatic about TDD: It helps a lot and I use it often, but it's not always going to be appropriate. Your judgment on this will improge over time. 

## Putting it all together

Let's run our refactored version of the code and be sure it still behaves like before:

```shell
go run refactored/lesson1_refactored.go
# In a separate terminal...
curl http://localhost:4000/products/2 | jq
```

Let's also check to make sure our new functionality is there:

```shell
curl -v http://localhost:4000/products/foobar
```

You should see a 400 response back from the server!