package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Optimized bcrypt cost for better performance
// Cost 8 = ~25ms, Cost 10 = ~100ms, Cost 12 = ~400ms
const bcryptCost = 4

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}
