package httptesting

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type HttpRequest struct {
	Method      string
	Path        string
	Body        interface{}
	Payload     string
	Description string
	Headers     map[string]string
}
type WriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w WriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

var file *os.File

func Prepare(filename string) {
	// Don't foget to call r.Use(MarkdownDebugLogger())

	if len(strings.TrimSpace(filename)) > 0 {
		var err error
		file, err = os.Create(filename)
		if err != nil {
			fmt.Errorf("cannot open %s: %v", filename, err)
		}
	}
}

func Teardown() {
	if file != nil {
		err := file.Close()
		if err != nil {
			fmt.Errorf("cannot close markdown file: %v", err)
		}
	}
}

func RegisterMarkdownDebugLogger(r *gin.Engine) {
	if len(r.Routes()) > 0 {
		fmt.Printf("ERROR: RegisterMarkdownDebugLogger() should be called before any other routes are registered\n")
	}
	r.Use(MarkdownDebugLogger())
}

func MarkdownDebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		description := c.Request.Header["__httptesting_desc"][0]
		if file != nil && len(description) > 0 {
			url := c.Request.URL.String()
			for _, p := range c.Params {
				url = strings.Replace(url, p.Value, ":"+p.Key, 1)
			}

			file.WriteString(fmt.Sprintf("\n* %s `%s` %s\n\n", c.Request.Method, url, description))
			for k, v := range c.Request.Header {
				for _, v1 := range v {
					if "__httptesting_desc" != k {
						file.WriteString(fmt.Sprintf("   - Header: `%s`: `%s`\n", k, v1))
					}
				}
			}
		}

		if c.Request.Body != nil {
			buf, _ := ioutil.ReadAll(c.Request.Body)
			reader1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			reader2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			requestBody := parseBody(reader1)

			if file != nil && len(description) > 0 {
				file.WriteString(fmt.Sprintf("\n   Request:\n```json\n%s\n```\n", requestBody))
			}

			c.Request.Body = reader2
		}

		wr := &WriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = wr
		c.Next()

		var response map[string]interface{}
		body := wr.Body.String()

		err := json.Unmarshal([]byte(body), &response)
		if err != nil {
			fmt.Errorf("Unable to parse the json response %d\n", err)
			if file != nil && len(description) > 0 {
				file.WriteString(fmt.Sprintf("\n   Response (%d):\n```text\n%s\n```\n", c.Writer.Status(), body))
			}
		} else {
			jsonDoc, _ := json.MarshalIndent(response, "", "\t")
			if file != nil && len(description) > 0 {
				file.WriteString(fmt.Sprintf("\n   Response (%d):\n```json\n%s\n```\n", c.Writer.Status(), jsonDoc))
			}
		}

	}
}

func parseBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

func PerformRequest(r *gin.Engine, request HttpRequest) *httptest.ResponseRecorder {
	var body io.Reader = nil
	if "GET" != request.Method {
		if request.Body != nil {
			jsonDoc, err := json.MarshalIndent(request.Body, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			body = bytes.NewBuffer(jsonDoc)

		} else if len(request.Payload) > 0 {
			body = bytes.NewBuffer([]byte(request.Payload))
		}
	}

	req, _ := http.NewRequest(request.Method, request.Path, body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("__httptesting_desc", request.Description)
	if request.Headers != nil {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func AssertStatusCode(t *testing.T, w *httptest.ResponseRecorder, expectedStatusCode int) {
	if w.Code != expectedStatusCode {
		t.Errorf("Unexpected status code: %d\n", w.Code)
	}
}

func AssertResponseStatus(t *testing.T, w *httptest.ResponseRecorder, expectedStatus string) map[string]interface{} {
	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("Unable to parse the json response %d: %s\n", err, w.Body.String())
		t.FailNow()
	}

	if response["Status"] != expectedStatus {
		fmt.Printf("Response: %s\n", w.Body.String())
		errorMessage := fmt.Sprintf("Unexpected status: %s, should be %s\n", response["Status"], expectedStatus)
		t.Error(errorMessage)
		panic(errors.New(errorMessage))
	}
	return response
}
