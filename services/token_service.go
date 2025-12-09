package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my-secret-key")

func CreateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,                                //storing the user id in the token
		"exp":     time.Now().Add(time.Hour * 24).Unix(), //token expires after 24 hours

	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
func ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["user_id"].(string)
		return userID, nil
	}

	return "", err
}
