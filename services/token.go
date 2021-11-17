package services

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const AK = "askldjfh321oim1249smnc"

// CustomClaims
type CustomClaims struct {
	App string
	jwt.StandardClaims
}

// ParseToken 解析 token
func ParseToken(data string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(data, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AK), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的token")
}
