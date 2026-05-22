package auth

import (
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"realtime-chat/src/config"
	"realtime-chat/src/database"
	"realtime-chat/src/models"
)

// Validation constants
const (
	MinUsernameLength = 3
	MaxUsernameLength = 30
	MinPasswordLength = 8
	MaxPasswordLength = 100
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateUsername validates username format and length
func ValidateUsername(username string) (string, error) {
	username = strings.TrimSpace(username)
	
	if len(username) < MinUsernameLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username must be at least 3 characters long")
	}
	if len(username) > MaxUsernameLength {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username must not exceed 30 characters")
	}
	
	if !usernameRegex.MatchString(username) {
		return "", fiber.NewError(fiber.StatusBadRequest, 
			"Username can only contain letters, numbers, underscores, and hyphens")
	}
	
	return html.EscapeString(username), nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < MinPasswordLength {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must be at least 8 characters long")
	}
	if len(password) > MaxPasswordLength {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must not exceed 100 characters")
	}
	
	hasUpper, hasLower, hasDigit := false, false, false
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
	
	if !hasUpper || !hasLower || !hasDigit {
		return fiber.NewError(fiber.StatusBadRequest, 
			"Password must contain at least one uppercase letter, one lowercase letter, and one number")
	}
	
	return nil
}

// SanitizeInput sanitizes general input
func SanitizeInput(input string, maxLength int) string {
	input = strings.TrimSpace(input)
	if len(input) > maxLength {
		input = input[:maxLength]
	}
	return html.EscapeString(input)
}

func SignUp(ctx *fiber.Ctx) error {
	var request models.AuthRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Println("Error parsing request body for Signup:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request format. Please check your data.",
		})
	}

	if request.Username == "" || request.Password == "" {
		log.Println("Username or password is empty")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username and password are required fields.",
		})
	}

	// Validate and sanitize username
	validatedUsername, err := ValidateUsername(request.Username)
	if err != nil {
		log.Println("Username validation failed:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	request.Username = validatedUsername

	// Validate password strength
	if err := ValidatePassword(request.Password); err != nil {
		log.Println("Password validation failed:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	var userExists database.DBUser
	if err := database.DB.Where("name = ?", request.Username).First(&userExists).Error; err == nil {
		log.Println("User already exists")
		return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error":   true,
			"message": "Username already exists. Please choose a different one.",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "An error occurred while processing your request. Please try again later.",
		})
	}

	user := database.DBUser{
		Name:     request.Username,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Error adding user to database: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "An error occurred while processing your request. Please try again later.",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Signup success",
	})
}

func Login(ctx *fiber.Ctx) error {
	var request models.AuthRequest
	if err := ctx.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body for Login: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request format. Please check your data.",
		})
	}

	if request.Username == "" || request.Password == "" {
		log.Println("Username or password is empty")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Username and password are required fields.",
		})
	}

	// Sanitize username (no strict validation on login)
	request.Username = SanitizeInput(request.Username, MaxUsernameLength)

	var user database.DBUser
	if err := database.DB.Where("name = ?", request.Username).First(&user).Error; err != nil {
		log.Printf("Error finding user in database: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid username or password.",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Printf("Password mismatch: %v", err)
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid username or password.",
		})
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	token, err := GenerateJWT(userID, user.Name)
	if err != nil {
		log.Printf("Error generating JWT: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "An error occurred while processing your request. Please try again later.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Login success",
		"token":   token,
	})
}

func GenerateJWT(userID string, username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	return token.SignedString([]byte(config.Config.JwtSecret))
}
