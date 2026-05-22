package auth

import (
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/config"
)

func AuthorizationMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	log.Printf("[Auth] Authorization header: %s", authHeader)
	
	if authHeader == "" {
		log.Println("[Auth] No Authorization header provided")
		return fiber.ErrUnauthorized
	}

	if err := ValidateJWTToken(authHeader); err != nil {
		log.Printf("[Auth] Error validating JWT token: %v", err)
		return fiber.ErrUnauthorized
	}

	// Parse token and set user info in context
	userID, userName := ParseJWTToken(authHeader)
	log.Printf("[Auth] Setting context - userID: %s, userName: %s", userID, userName)
	
	ctx.Locals("userID", userID)
	ctx.Locals("userName", userName)

	return ctx.Next()
}

func ValidateJWTToken(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JwtSecret), nil
	})

	return err
}

func ParseJWTToken(tokenString string) (userID string, userName string) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID = claims["userID"].(string)
		userName = claims["username"].(string)
	}

	return userID, userName
}
