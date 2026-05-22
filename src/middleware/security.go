package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent clickjacking attacks
		c.Set("X-Frame-Options", "DENY")
		
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")
		
		// Enable XSS protection
		c.Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy (adjust based on your needs)
		c.Set("Content-Security-Policy", 
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self' data:; "+
			"connect-src 'self' ws: wss:;")
		
		// Strict Transport Security (HTTPS only - enable in production)
		// Uncomment when using HTTPS
		// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		// Permissions Policy (formerly Feature Policy)
		c.Set("Permissions-Policy", 
			"geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}

// CORS configuration with environment-based origins
func ConfigureCORS(allowedOrigins []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Set("Access-Control-Allow-Origin", origin)
		}

		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
