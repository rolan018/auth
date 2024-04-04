package jwt

import (
	"auth/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrExpiredToken = errors.New("token was expired")
)

func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string, app models.App) error {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})
	if err != nil {
		return err
	}
	expirationUnixTime := claims["exp"].(int64)
	if expirationUnixTime < time.Now().Unix() {
		return ErrExpiredToken
	}
	return nil
}
