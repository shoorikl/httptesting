package httptesting

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type RequestType struct {
	Status string `json:"Status"`
}

func TestMain(m *testing.M) {
	Prepare("chitchat.md")
	code := m.Run()
	Teardown()
	os.Exit(code)
}

func TestGETRequest(t *testing.T) {
	r := createRouter()
	w := PerformRequest(r, HttpRequest{Method: "GET", Path: "/test", Description: "Test GET Endpoint"})
	UnwrapResponse(t, w, "OK")
}

func TestPOSTRequest(t *testing.T) {
	req := gin.H{"Status": "HELLO"}
	r := createRouter()
	w := PerformRequest(r, HttpRequest{Method: "POST", Path: "/echo", Description: "Test POST Endpoint", Body: req})
	UnwrapResponse(t, w, "HELLO")
}

func createRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status": "OK"})
	})

	r.POST("/echo", func(c *gin.Context) {
		req := RequestType{}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(500, gin.H{
				"Status": "Error",
				"Error":  fmt.Sprintf("%v", err)})

		} else {
			c.JSON(200, gin.H{
				"Status": req.Status})
		}
	})
	return r
}
