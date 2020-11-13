package httptesting

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func getJwtSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) == 0 {
		jwtSecret = "jwtsecret"
	}
	return jwtSecret
}

// CreateToken create JWT token given the payload, expiration and authorization flag
func CreateToken(payload string, duration time.Duration, authorized bool) (string, error) {
	var err error

	claim := jwt.MapClaims{}
	claim["authorized"] = authorized
	claim["payload"] = payload
	claim["exp"] = time.Now().Add(duration).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := at.SignedString([]byte(getJwtSecret()))
	if err != nil {
		return "", err
	}
	return token, nil
}

func verifyToken(token string) (*jwt.Token, error) {

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Error: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(getJwtSecret()), nil
	})
	if err != nil {
		return nil, err
	}
	return jwtToken, nil
}

func verifyAuthorizationToken(token string) (valid bool, claims jwt.MapClaims) {
	jwtToken, err := verifyToken(token)

	if err != nil {
		return false, nil
	}

	if err != nil {
		fmt.Printf("Error: Error validating JWT token: %s\n", err.Error())
		return false, nil
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Error: No valid calaims found in JWT token\n")
		return false, nil
	}

	if !jwtToken.Valid {
		fmt.Printf("Error: JWT Token is not valid\n")
		return false, claims
	}

	return true, claims
}

// VerifyAuthorization validate authorization header
func VerifyAuthorization(c *gin.Context) (valid bool, token string, claims jwt.MapClaims) {
	if len(c.Request.Header["Authorization"]) > 0 {
		authorizationHeader := c.Request.Header["Authorization"][0]
		token := strings.Replace(authorizationHeader, "Bearer ", "", 1)
		valid, claims := verifyAuthorizationToken(token)
		if !valid {
			return false, token, claims
		} else {
			if false == claims["authorized"].(bool) {
				return false, token, claims
			} else {
				return true, token, claims
			}
		}
	}
	return false, "", nil
}

// ExtractPayloadField extract a field from a json JWT claim payload
func ExtractPayloadField(claims jwt.MapClaims, field string) string {
	if claims["payload"] == nil {
		return ""
	}
	payloadJSON := claims["payload"].(string)
	var payload map[string]interface{}
	json.Unmarshal([]byte(payloadJSON), &payload)
	value := payload[field].(string)
	return value
}

// CreatePayloadField create a JSON JWT claim payload from a map
func CreatePayloadField(payload map[string]interface{}) string {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error converting to JSON: %s\n", err.Error())
		return "{}"
	}
	return string(jsonPayload)
}

// ExtractPayloadNumericField extract a field from a json JWT claim payload
func ExtractPayloadNumericField(claims jwt.MapClaims, field string) float64 {
	if claims["payload"] == nil {
		return 0.0
	}
	payloadJSON := claims["payload"].(string)
	var payload map[string]interface{}
	json.Unmarshal([]byte(payloadJSON), &payload)
	value := payload[field].(float64)
	return value
}
