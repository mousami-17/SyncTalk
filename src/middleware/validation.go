package middleware

import (
	"html"
	"regexp"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// Validation rules
const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
	MinPasswordLength = 8
	MaxPasswordLength = 100
	MaxMessageLength  = 2000
	MaxRoomNameLength = 50
)

// Username validation regex (alphanumeric, underscore, hyphen)
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateUsername validates username format and length
func ValidateUsername(username string) (string, error) {
	// Trim whitespace
	username = strings.TrimSpace(username)

	// Check length
	if len(username) < MinUsernameLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username must be at least 3 characters long")
	}
	if len(username) > MaxUsernameLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username must not exceed 30 characters")
	}

	// Check format (alphanumeric, underscore, hyphen only)
	if !usernameRegex.MatchString(username) {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username can only contain letters, numbers, underscores, and hyphens")
	}

	// Sanitize (escape HTML)
	username = html.EscapeString(username)

	return username, nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	// Check length
	if len(password) < MinPasswordLength {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must be at least 8 characters long")
	}
	if len(password) > MaxPasswordLength {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must not exceed 100 characters")
	}

	// Check for at least one uppercase letter
	hasUpper := false
	// Check for at least one lowercase letter
	hasLower := false
	// Check for at least one digit
	hasDigit := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if !hasUpper {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must contain at least one number")
	}

	return nil
}

// SanitizeMessage sanitizes message content
func SanitizeMessage(message string) (string, error) {
	// Trim whitespace
	message = strings.TrimSpace(message)

	// Check if empty
	if message == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Message cannot be empty")
	}

	// Check length
	if len(message) > MaxMessageLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Message must not exceed 2000 characters")
	}

	// Escape HTML to prevent XSS
	message = html.EscapeString(message)

	return message, nil
}

// ValidateRoomName validates room name
func ValidateRoomName(roomName string) (string, error) {
	// Trim whitespace
	roomName = strings.TrimSpace(roomName)

	// Check if empty
	if roomName == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Room name cannot be empty")
	}

	// Check length
	if len(roomName) > MaxRoomNameLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Room name must not exceed 50 characters")
	}

	// Sanitize (escape HTML)
	roomName = html.EscapeString(roomName)

	return roomName, nil
}

// SanitizeInput is a general purpose input sanitizer
func SanitizeInput(input string, maxLength int) string {
	// Trim whitespace
	input = strings.TrimSpace(input)

	// Truncate if too long
	if len(input) > maxLength {
		input = input[:maxLength]
	}

	// Escape HTML
	input = html.EscapeString(input)

	return input
}
