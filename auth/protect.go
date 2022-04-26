package auth

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

func Protect(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("unexpected singing method: %v", token.Header["alg"])
		}

		signingKey := []byte("==signature==")
		return signingKey, nil
	})

	return err
}
