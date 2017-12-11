package services

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func generateJwt(accountId string) (tokenString string, err error) {
	expiresTime := time.Now().Add(time.Hour)
	userToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   accountId,
		ExpiresAt: expiresTime.Unix(),
	})
	tokenString, err = userToken.SignedString([]byte("0123456789"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
