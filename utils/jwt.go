// this package creates a jwt for login
package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getting secret key from .env using the init function
var secretKey []byte

func init() {

	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		log.Fatal("SECRET_KEY not found in environment")
	}
	secretKey = []byte(secret_key)
}

func CreateToken(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// validating the token
func ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &claims, nil
	}
	return nil, fmt.Errorf("Inavlid token claims")
}
