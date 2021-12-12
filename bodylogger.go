package httptesting

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type RequestLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w RequestLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var whitelistedIps = getWhitelistedIps()

func getWhitelistedIps() map[string]bool {
	whitelistedIps := make(map[string]bool)

	whitelist := os.Getenv("SWAGGER_WHITELIST")
	if len(whitelist) == 0 {
		whitelistedIps["*"] = true
	} else {
		ips := strings.Split(whitelist, ",")
		for _, ip := range ips {
			whitelistedIps[ip] = true
		}
	}
	return whitelistedIps
}

func SwaggerLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.Contains(c.Request.URL.RequestURI(), "/swagger.json") {
			remoteAddr, _, _ := net.SplitHostPort(c.Request.RemoteAddr)

			fmt.Printf("Swagger url %s accessed from %s\n", c.Request.URL, remoteAddr)

			if _, ok := whitelistedIps["*"]; !ok {
				if _, ok := whitelistedIps[remoteAddr]; !ok {
					c.AbortWithError(http.StatusUnauthorized, errors.New("You're not allowed to access this endpoint"))
					return
				}
			}
		}
	}
}

func BodyLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		health := strings.Contains(c.Request.URL.RequestURI(), "/healthz")
		routeDiscovery := strings.Contains(c.Request.URL.RequestURI(), "/routes")
		graphQlPlayground := strings.Contains(c.Request.URL.RequestURI(), "/graphql")
		graphQlQuery := strings.Contains(c.Request.URL.RequestURI(), "/query")
		rootQuery := c.Request.URL.RequestURI() == "/"

		logRequest := !health && !routeDiscovery && !graphQlPlayground && !graphQlQuery && !rootQuery
		if logRequest {
			if "GET" == c.Request.Method {
				fmt.Printf("\nRequest: %s %s\n", c.Request.Method, c.Request.URL.RequestURI())
			} else {
				var body []byte
				if c.Request.Body != nil {
					body, _ = ioutil.ReadAll(c.Request.Body)
				}
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				fmt.Printf("\nRequest: %s %s Body: %s\n", c.Request.Method, c.Request.URL.RequestURI(), string(body))
			}
			for name, values := range c.Request.Header {
				if strings.Index(name, "__httptesting") != 0 {
					for _, value := range values {
						fmt.Printf("  * %s=%s\n", name, value)
					}
				}
			}
		}

		blw := &RequestLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		if logRequest {
			contentType, found := blw.Header()["Content-Type"]
			if found && len(contentType) > 0 {
				if strings.Contains(contentType[0], "application/json") {
					fmt.Printf("Response: [%d] Body: %s\n", c.Writer.Status(), blw.body.String())

					for name, values := range blw.Header() {
						for _, value := range values {
							fmt.Printf("  * %s=%s\n", name, value)
						}
					}

				} else {
					fmt.Printf("Non-json response\n")
				}
			}
		}
	}
}
