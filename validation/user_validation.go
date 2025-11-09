package validation

import (
	"errors"
	"net/mail"
	"regexp"
	"unicode/utf8"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 20
)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("Email cannot be empty")
	}
	//parsing the email
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("Invalid email format")
	}

	return nil
}
func ValidatePassword(password string) error {
	passwordLength := utf8.RuneCountInString(password)

	if passwordLength < minPasswordLength || passwordLength > maxPasswordLength {
		return errors.New("password must be between 8 and 20 characters long")
	}
	// Add these checks for stronger passwords:
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain uppercase, lowercase, and number")
	}
	return nil
}

// validate username
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username)
	if !validUsername {
		return errors.New("username can only contain letters, numbers, and underscores")
	}
	return nil
}
