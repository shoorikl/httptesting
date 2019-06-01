# httptesting

This package somewhat simplifies http endpoint testing

# Installing
`go get -u github.com/shoorikl/httptesting`
`dep ensure -add github.com/shoorikl/httptesting`

In you test package, add

```go
func TestMain(m *testing.M) {

	httptesting.Prepare()
	code := m.Run()
	httptesting.Teardown()
	os.Exit(code)
}
```

# Usage

```
var req map[string]interface{}
...
w := httptesting.PerformRequestWithBody(r, "POST", "/some/endpoint", req)
response := httptesting.UnwrapResponse(t, w, "OK")
```



