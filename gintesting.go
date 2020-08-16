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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoisie/mustache"
)

type HttpRequest struct {
	Method            string
	Path              string
	Body              interface{}
	Payload           string
	Description       string
	Headers           map[string]string
	ResponseVariables []ResponseVariable
	Name              string
}

type ResponseVariable struct {
	Variable   string
	Expression string
	Value      interface{}
}

type WriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w WriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

var docFile *os.File
var httpFile *os.File
var baseUrl string
var variables map[string]interface{} = make(map[string]interface{})

func Prepare(docFileName string) {
	// Don't foget to call r.Use(MarkdownDebugLogger())

	if len(strings.TrimSpace(docFileName)) > 0 {
		var err error
		docFile, err = os.Create(docFileName)
		if err != nil {
			fmt.Errorf("cannot open %s: %v", docFileName, err)
		}
	}
}

func PrepareWithHttpDoc(docFileName string, httpFileName string, baseUrlParam string) {
	// Don't foget to call r.Use(MarkdownDebugLogger())

	if len(strings.TrimSpace(docFileName)) > 0 {
		var err error
		docFile, err = os.Create(docFileName)
		if err != nil {
			fmt.Errorf("cannot open %s: %v", docFileName, err)
		}
	}

	if len(strings.TrimSpace(httpFileName)) > 0 {
		var err error
		httpFile, err = os.Create(httpFileName)
		if err != nil {
			fmt.Errorf("cannot open %s: %v", httpFileName, err)
		} else {
			httpFile.WriteString(fmt.Sprintf("@baseUrl = %s\n\n", baseUrlParam))
		}
	}
	baseUrl = baseUrlParam
}

func Teardown() {
	if docFile != nil {
		err := docFile.Close()
		if err != nil {
			fmt.Errorf("cannot close markdown file: %v", err)
		}
	}

	if httpFile != nil {
		err := httpFile.Close()
		if err != nil {
			fmt.Errorf("cannot close http file: %v", err)
		}
	}
}

func RegisterMarkdownDebugLogger(r *gin.Engine) {
	if len(r.Routes()) > 0 {
		fmt.Printf("ERROR: RegisterMarkdownDebugLogger() should be called before any other routes are registered\n")
	}
	r.Use(MarkdownDebugLogger())
}

func ExtractVariables(w *httptest.ResponseRecorder, responseVariables []ResponseVariable) {

	if len(responseVariables) > 0 {
		httpFile.WriteString("###\n")
		httpFile.WriteString("\n")
		for _, responseVariable := range responseVariables {
			httpFile.WriteString(fmt.Sprintf("@%s = {{%s}}\n", responseVariable.Variable, responseVariable.Expression))
			variables[responseVariable.Variable] = responseVariable.Value
		}
		httpFile.WriteString("\n")
	}
}

func MarkdownDebugLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		description := c.Request.Header["__httptesting_desc"][0]
		name := c.Request.Header["__httptesting_name"][0]

		if docFile != nil && len(description) > 0 {
			url := c.Request.URL.String()
			for _, p := range c.Params {
				url = strings.Replace(url, p.Value, ":"+p.Key, 1)
			}

			if httpFile != nil {
				httpFile.WriteString("###\n")

				httpFile.WriteString(fmt.Sprintf("# %s\n", description))
				if len(name) > 0 {
					httpFile.WriteString(fmt.Sprintf("# @name %s\n", name))
				}
				httpFile.WriteString(fmt.Sprintf("%s {{baseUrl}}%s\n", c.Request.Method, url))
			}

			docFile.WriteString(fmt.Sprintf("\n* %s `%s` %s\n\n", c.Request.Method, url, description))
			docFile.WriteString("   - Request:\n")
			if len(c.Request.Header) > 0 {
				docFile.WriteString("      - Headers:\n")

				for k, v := range c.Request.Header {
					for i, v1 := range v {
						if strings.Index(k, "__httptesting") != 0 {

							if httpFile != nil {
								httpFile.WriteString(fmt.Sprintf("%s: %s\n", k, v1))
							}
							v[i] = PopulateVariables(v1)
							docFile.WriteString(fmt.Sprintf("         - `%s`: `%s`\n", k, v1))
						}

					}
				}
				httpFile.WriteString("\n")
			}

			if c.Request.Body != nil {
				buf, _ := ioutil.ReadAll(c.Request.Body)
				reader1 := ioutil.NopCloser(bytes.NewBuffer(buf))
				reader2 := ioutil.NopCloser(bytes.NewBuffer(buf))
				body := parseBody(reader1)
				requestBody := indent(body)

				docFile.WriteString(fmt.Sprintf("      - Body:\n\t\t```json\n%s\t\t```\n", requestBody))
				if httpFile != nil {
					httpFile.WriteString(fmt.Sprintf("%s\n", body))
				}

				c.Request.Body = reader2
			}

			if httpFile != nil {
				httpFile.WriteString("\n")
			}
		}

		wr := &WriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = wr
		c.Next()

		var response map[string]interface{}
		body := wr.Body.String()

		docFile.WriteString(fmt.Sprintf("\n   - Response (%d)\n", c.Writer.Status()))

		if len(c.Writer.Header()) > 0 {
			docFile.WriteString("      - Headers:\n")

			for k, v := range c.Writer.Header() {
				for _, v1 := range v {
					docFile.WriteString(fmt.Sprintf("         - `%s`: `%s`\n", k, v1))
				}
			}
		}

		err := json.Unmarshal([]byte(body), &response)
		if err != nil {
			fmt.Errorf("Unable to parse the json response %d\n", err)
			if docFile != nil && len(description) > 0 {
				docFile.WriteString(fmt.Sprintf("\n      - Body:\n\t\t```text\n%s\t\t```\n", indent(body)))
			}
		} else {
			jsonDoc, _ := json.MarshalIndent(response, "", "\t")

			if docFile != nil && len(description) > 0 {
				docFile.WriteString(fmt.Sprintf("\n      - Body:\n\t\t```json\n%s\t\t```\n", indent(string(jsonDoc))))
			}
		}

	}
}

func PopulateVariables(template string) string {
	return mustache.Render(template, variables)
}

func parseBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	return buf.String()
}

func indent(body string) string {
	lines := strings.Split(body, "\n")
	sb := StringBuilder{}
	for _, line := range lines {
		if len(strings.TrimSpace(line)) > 0 {
			sb.Write("\t\t").Write(line).Write("\n")
		}
	}
	return sb.String()
}

// Makes a call to a url exposed by a Gin engine, logging request and a response
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
	req.Header.Set("__httptesting_name", request.Name)

	if request.Headers != nil {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

// Makes a call to a fully qualified remote url
func PerformRemoteRequest(request HttpRequest) ([]byte, error) {
	var body io.Reader = nil
	if "GET" != request.Method {
		if request.Body != nil {
			jsonDoc, err := json.MarshalIndent(request.Body, "", "\t")
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				return nil, err
			}
			body = bytes.NewBuffer(jsonDoc)

		} else if len(request.Payload) > 0 {
			body = bytes.NewBuffer([]byte(request.Payload))
		}
	}

	req, _ := http.NewRequest(request.Method, request.Path, body)
	req.Header.Set("Content-Type", "application/json")

	if request.Headers != nil {
		for k, v := range request.Headers {
			req.Header.Set(k, v)
		}
	}

	client := http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return nil, err
	}

	return bodyBytes, nil
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization, content-type, accept, origin, Cache-Control, X-Requested-With, access-control-allow-origin, access-control-allow-credentials, access-control-allow-headers, access-control-allow-methods")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if "OPTIONS" == c.Request.Method {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
