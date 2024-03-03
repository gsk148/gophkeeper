package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secret = []byte("adfw5F#s8deRefyd")

func EncodeToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"sub": uid,
	})

	return token.SignedString(secret)
}

func GetUserIDFromToken(token string) (string, error) {
	claims, err := getClaims(token)
	if err != nil {
		return "", err
	}
	return claims["sub"].(string), nil
}

func IsTokenExpired(token string) (bool, error) {
	claims, err := getClaims(token)
	if err != nil {
		return true, err
	}
	return !claims.VerifyExpiresAt(jwt.TimeFunc().Unix(), true), nil
}

func getClaims(ts string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected token signing method")
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return jwt.MapClaims{}, errors.New("couldn't get claims")
}
