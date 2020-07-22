package httptesting

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

func BodyLogger(c *gin.Context) {
	health := strings.Contains(c.Request.URL.RequestURI(), "/healthz")
	if !health {
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
	}

	blw := &RequestLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()

	if !health {
		fmt.Printf("Response: [%d] Body: %s\n", c.Writer.Status(), blw.body.String())
	}
}
