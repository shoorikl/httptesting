# httptesting

Restful API markdown logging middleware for Golang Gin http routing framework.

# Conventions

By design, every endpoint has a mandatory `Status` response filed, which conveys granular details about the outcome of a call. In case of an error, in addition of HTTP 500, a field called `Error` is expected to be returned. In case of a success, in addition to the HTTP 200, `Status` could contian information like `Inserted`, `Deleted`, etc. This is ultimately a semantic confirmation.

# Installing
`go get -u github.com/shoorikl/httptesting`
`dep ensure -add github.com/shoorikl/httptesting`

In you test package, add

```go
var r = createRouter()

func TestMain(m *testing.M) {

	httptesting.Prepare("chitchat.md")
	code := m.Run()
	httptesting.Teardown()
	os.Exit(code)
}

func createRouter() *gin.Engine {
	r := gin.Default()
	httptesting.RegisterMarkdownDebugLogger(r)
	return r
}
```

# Usage

```
	req := gin.H{"Status": "HELLO"}
	r := createRouter() // initialize gin.Engine
	...
	w := httptesting.PerformRequest(r, httptesting.HttpRequest{Method: "POST", Path: "/echo", Description: "Test POST Endpoint", Body: req})
	AssertResponseStatus(t, w, "HELLO")
```

This will execute a POST call to /echo, and assert the Status field of the response payload. When running `go test`, a `chitchat.md` file will be created with all request-response examples.

# Outcome

By running `go test` on a test package that is instrumented with httptest, you will receive a markdown snippet of all interactions with your endpoints, including http methods, uris, request and response payloads -- which is useful in addition to the OpenAPI/Swagger, as it gives you concrete examples.

Below is a sample from this project:



* GET `/test` Test GET Endpoint

   Header: `Content-Type`: `application/json`

   Response (200):
```json
{
	"Status": "OK"
}
```

* POST `/echo` Test POST Endpoint

   Header: `Content-Type`: `application/json`
   Header: `Token`: `123`

   Request:
```json
{
	"Status": "HELLO"
}
```

   Response (200):
```json
{
	"Status": "HELLO"
}
```

* GET `/param/:value` Test GET Endpoint with route param

   Header: `Content-Type`: `application/json`

   Response (200):
```json
{
	"Status": "somevalue"
}
```

* PUT `/param/:value` Test PUT Endpoint with route param

   Header: `Content-Type`: `application/json`

   Request:
```json
{
	"Status": "HELLO"
}
```

   Response (200):
```json
{
	"Status": "somevalue"
}
```
