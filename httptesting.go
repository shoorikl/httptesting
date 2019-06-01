package httptesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type HttpRequest struct {
	Method      string
	Path        string
	Body        interface{}
	Payload     string
	Description string
}

var file *os.File

func Prepare(filename string) {
	if len(strings.TrimSpace(filename)) > 0 {
		var err error
		file, err = os.Create(filename)
		if err != nil {
			fmt.Errorf("Cannot open chitchat.md: %v\n", err)
		}
	}
}

func Teardown() {
	if file != nil {
		file.Close()
	}
}

func PerformRequest(r http.Handler, request HttpRequest) *httptest.ResponseRecorder {
	if file != nil {
		file.WriteString(fmt.Sprintf("\n* %s `%s` %s\n", request.Method, request.Path, request.Description))
	}

	var body io.Reader = nil
	if "GET" != request.Method {
		if request.Body != nil {
			jsonDoc, err := json.MarshalIndent(request.Body, "", "\t")
			if err != nil {
				log.Fatal(err)
			}
			body = bytes.NewBuffer(jsonDoc)
			fmt.Printf("Request: %s\n", jsonDoc)
			if file != nil {
				file.WriteString(fmt.Sprintf("\nRequest:\n```json\n%s\n```\n", jsonDoc))
			}
		} else if len(request.Payload) > 0 {
			body = bytes.NewBuffer([]byte(request.Payload))
			fmt.Printf("Request: %s\n", request.Payload)
			if file != nil {
				file.WriteString(fmt.Sprintf("\nRequest:\n```json\n%s\n```\n", request.Payload))
			}
		}
	}

	req, _ := http.NewRequest(request.Method, request.Path, body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func UnwrapResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus string) map[string]interface{} {
	var response map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Errorf("Unable to parse the json response %d\n", err)
		t.FailNow()
	}
	jsonDoc, _ := json.MarshalIndent(response, "", "\t")
	if file != nil {
		file.WriteString(fmt.Sprintf("\nResponse:\n```json\n%s\n```\n", jsonDoc))
	}
	fmt.Printf("Response: %s\n", jsonDoc)
	if response["Status"] != expectedStatus {
		t.Errorf("Unexpected status: %s\n", response["Status"])
	}
	return response
}
