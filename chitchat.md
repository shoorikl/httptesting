
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
