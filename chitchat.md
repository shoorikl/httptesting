
* GET `/test` Test GET Endpoint

Response (200):
```json
{
	"Status": "OK"
}
```

* POST `/echo` Test POST Endpoint

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

Response (200):
```json
{
	"Status": "somevalue"
}
```

* PUT `/param/:value` Test PUT Endpoint with route param

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
