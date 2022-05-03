package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// middleware
func Protect(signature []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorization, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
			}

			return signature, nil
		})

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized) // stop everything
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok {
			aud := claims["aud"] // audian
			c.Set("aud", aud)
		}

		c.Next() // next middleware
	}
}
