package gintesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var file *os.File

func prepare() {
	var err error
	file, err = os.Create("chitchat.md")
	if err != nil {
		fmt.Errorf("Cannot open chitchat.md: %v\n", err)
	}
}

func teardown() {
	if file != nil {
		file.Close()
	}
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("Request: []\n")
	if file != nil {
		file.WriteString(fmt.Sprintf("\n* %s `%s`\n", method, path))
	}
	return w
}

func performRequestWithBody(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	jsonDoc, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	if file != nil {
		file.WriteString(fmt.Sprintf("\n* %s `%s`\n", method, path))
		file.WriteString(fmt.Sprintf("\nRequest:\n```json\n%s\n```\n", jsonDoc))
	}

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonDoc))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	fmt.Printf("Request: %s\n", jsonDoc)
	return w
}

func unwrapResponse(t *testing.T, w *httptest.ResponseRecorder, expectedStatus string) map[string]interface{} {
	var response map[string]interface{}
	fmt.Printf("Response: %s\n", w.Body.String())

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
