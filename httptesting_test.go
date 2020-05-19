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

var r = createRouter()

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
	w := PerformRequest(r, HttpRequest{Method: "POST", Path: "/echo", Description: "Test POST Endpoint", Body: req, Headers: map[string]string{"Token": "123"}})
	AssertResponseStatus(t, w, "HELLO")
}

func TestGETRouteParamRequest(t *testing.T) {

	w := PerformRequest(r, HttpRequest{Method: "GET", Path: "/param/somevalue", Description: "Test GET Endpoint with route param"})
	AssertResponseStatus(t, w, "somevalue")
}

func TestPUTRouteParamRequest(t *testing.T) {
	req := gin.H{"Status": "HELLO"}
	w := PerformRequest(r, HttpRequest{Method: "PUT", Path: "/param/somevalue", Description: "Test PUT Endpoint with route param", Body: req})
	AssertResponseStatus(t, w, "somevalue")
}

func createRouter() *gin.Engine {
	r := gin.Default()
	RegisterMarkdownDebugLogger(r)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status": "OK"})
	})

	r.GET("/param/:value", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"Status": c.Param("value")})
	})

	r.PUT("/param/:value", func(c *gin.Context) {
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

func TestSelectStatement(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Select("mytable").Returning("firstname", "lastname").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "SELECT firstname, lastname FROM mytable WHERE customerid=$1 AND accounttype=$2 AND active=$3" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 3 {
		t.Errorf("Should take 3 value arguments: %v", inArgs)
	}

	if len(outArgs) != 2 {
		t.Errorf("Should return 2 parameters: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestSelectStatementWithLimit(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Select("mytable").Returning("firstname", "lastname").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Limit(3).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "SELECT firstname, lastname FROM mytable WHERE customerid=$1 AND accounttype=$2 AND active=$3 LIMIT 3" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 3 {
		t.Errorf("Should take 3 value arguments: %v", inArgs)
	}

	if len(outArgs) != 2 {
		t.Errorf("Should return 2 parameters: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestSelectStatementError(t *testing.T) {
	sb := PostgresSqlBuilder{}
	_, _, _, err := sb.Select("").Returning("firstname", "lastname").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err == nil {
		t.Errorf("Should return an error")
	}

	_, _, _, err = sb.Select("mytable").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err == nil {
		t.Errorf("Should return an error")
	}

	_, _, _, err = sb.Select("mytable").Returning("firstname", "lastname").Build()

	if err == nil {
		t.Errorf("Should return an error")
	}

}

func TestInsertStatement(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Insert("mytable").Returning("id").SetArg("lifetimevalue", 100).SetArg("customerid", 5).SetArg("accounttype", "seller").SetArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "INSERT INTO mytable (lifetimevalue, customerid, accounttype, active) VALUES ($1, $2, $3, $4) RETURNING id" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 4 {
		t.Errorf("Should take 4 value arguments: %v", inArgs)
	}

	if len(outArgs) != 1 {
		t.Errorf("Should return 1 parameter: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestUpdateStatement(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Update("mytable").Returning("id").SetArg("lifetimevalue", 100).SetArg("active", false).WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "UPDATE mytable lifetimevalue=$1,active=$2 WHERE customerid=$3 AND accounttype=$4 AND active=$5 RETURNING id" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 5 {
		t.Errorf("Should take 5 value arguments: %v", inArgs)
	}

	if len(outArgs) != 1 {
		t.Errorf("Should return 1 parameter: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestUpdateStatementNoReturning(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Update("mytable").SetArg("lifetimevalue", 100).SetArg("active", false).WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "UPDATE mytable lifetimevalue=$1,active=$2 WHERE customerid=$3 AND accounttype=$4 AND active=$5" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 5 {
		t.Errorf("Should take 5 value arguments: %v", inArgs)
	}

	if len(outArgs) != 0 {
		t.Errorf("Should return 0 parameters: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestDeleteStatement(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Delete("mytable").Returning("id").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "DELETE FROM mytable WHERE customerid=$1 AND accounttype=$2 AND active=$3 RETURNING id" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 3 {
		t.Errorf("Should take 5 value arguments: %v", inArgs)
	}

	if len(outArgs) != 1 {
		t.Errorf("Should return 1 parameter: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}

func TestDeleteStatementNoReturning(t *testing.T) {
	sb := PostgresSqlBuilder{}
	query, inArgs, outArgs, err := sb.Delete("mytable").WhereArg("customerid", 5).WhereArg("accounttype", "seller").WhereArg("active", true).Build()

	if err != nil {
		t.Errorf("Cannot create query: %s", err.Error())
	}

	if query != "DELETE FROM mytable WHERE customerid=$1 AND accounttype=$2 AND active=$3" {
		t.Errorf("Mismatching query: %s\n", query)
	}

	if len(inArgs) != 3 {
		t.Errorf("Should take 5 value arguments: %v", inArgs)
	}

	if len(outArgs) != 0 {
		t.Errorf("Should return 0 parameters: %v", outArgs)
	}

	t.Logf("In Args: %v\n", inArgs)
}
