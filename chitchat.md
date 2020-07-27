
* GET `/test` Test GET Endpoint

   - Request:
      - Headers:
         - `Content-Type`: `application/json`

   - Response (200)
      - Headers:
         - `Access-Control-Allow-Headers`: `Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods`
         - `Access-Control-Allow-Methods`: `POST, OPTIONS, GET, PUT, DELETE`
         - `Content-Type`: `application/json; charset=utf-8`
         - `Access-Control-Allow-Origin`: `*`
         - `Access-Control-Allow-Credentials`: `true`

      - Body:
		```json
		{
			"Status": "OK"
		}
		```

* POST `/echo` Test POST Endpoint

   - Request:
      - Headers:
         - `Token`: `123`
         - `Content-Type`: `application/json`
      - Body:
		```json
		{
			"Status": "HELLO"
		}
		```

   - Response (200)
      - Headers:
         - `Access-Control-Allow-Origin`: `*`
         - `Access-Control-Allow-Credentials`: `true`
         - `Access-Control-Allow-Headers`: `Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods`
         - `Access-Control-Allow-Methods`: `POST, OPTIONS, GET, PUT, DELETE`
         - `Content-Type`: `application/json; charset=utf-8`

      - Body:
		```json
		{
			"Status": "HELLO"
		}
		```

* POST `/login` Test POST Auth Endpoint

   - Request:
      - Headers:
         - `Content-Type`: `application/json`
         - `Token`: `123`
      - Body:
		```json
		{
			"Status": "HELLO"
		}
		```

   - Response (200)
      - Headers:
         - `Content-Type`: `application/json; charset=utf-8`
         - `Access-Control-Allow-Origin`: `*`
         - `Access-Control-Allow-Credentials`: `true`
         - `Access-Control-Allow-Headers`: `Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods`
         - `Access-Control-Allow-Methods`: `POST, OPTIONS, GET, PUT, DELETE`

      - Body:
		```json
		{
			"AuthToken": "token body",
			"Status": "HELLO"
		}
		```

* GET `/param/:value` Test GET Endpoint with route param

   - Request:
      - Headers:
         - `Content-Type`: `application/json`

   - Response (200)
      - Headers:
         - `Access-Control-Allow-Origin`: `*`
         - `Access-Control-Allow-Credentials`: `true`
         - `Access-Control-Allow-Headers`: `Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods`
         - `Access-Control-Allow-Methods`: `POST, OPTIONS, GET, PUT, DELETE`
         - `Content-Type`: `application/json; charset=utf-8`

      - Body:
		```json
		{
			"Status": "somevalue"
		}
		```

* PUT `/param/:value` Test PUT Endpoint with route param

   - Request:
      - Headers:
         - `Content-Type`: `application/json`
      - Body:
		```json
		{
			"Status": "HELLO"
		}
		```

   - Response (200)
      - Headers:
         - `Access-Control-Allow-Origin`: `*`
         - `Access-Control-Allow-Credentials`: `true`
         - `Access-Control-Allow-Headers`: `Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods`
         - `Access-Control-Allow-Methods`: `POST, OPTIONS, GET, PUT, DELETE`
         - `Content-Type`: `application/json; charset=utf-8`

      - Body:
		```json
		{
			"Status": "somevalue"
		}
		```
