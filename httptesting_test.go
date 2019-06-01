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

var r *gin.Engine = createRouter()

func TestMain(m *testing.M) {

	Prepare("chitchat.md")
	code := m.Run()
	Teardown()
	os.Exit(code)
}

func TestGETRequest(t *testing.T) {

	w := PerformRequest(r, HttpRequest{Method: "GET", Path: "/test", Description: "Test GET Endpoint"})
	AssertResponseStatus(t, w, "OK")
}

func TestPOSTRequest(t *testing.T) {

	req := gin.H{"Status": "HELLO"}
	w := PerformRequest(r, HttpRequest{Method: "POST", Path: "/echo", Description: "Test POST Endpoint", Body: req})
	AssertResponseStatus(t, w, "HELLO")
}

func TestGETRouteParamRequest(t *testing.T) {

	w := PerformRequest(r, HttpRequest{Method: "GET", Path: "/param/somevalue", Description: "Test GET Endpoint with route param"})
	AssertResponseStatus(t, w, "somevalue")
}

func TestInvalidPUTRequest(t *testing.T) {

	w := PerformRequest(r, HttpRequest{Method: "PUT", Path: "/invalidpath", Description: "Test Invalid PUT request"})
	AssertResponseStatus(t, w, "OK")
}

func createRouter() *gin.Engine {
	r := gin.Default()
	r.Use(MarkdownDebugLogger())

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status": "OK"})
	})

	r.GET("/param/:value", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status": c.Param("value")})
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
