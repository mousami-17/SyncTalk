package main

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"realtime-chat/src/auth"
	"realtime-chat/src/cache"
	"realtime-chat/src/chat"
	"realtime-chat/src/config"
	"realtime-chat/src/database"
	"realtime-chat/src/middleware"
	"realtime-chat/src/upload"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler:          customErrorHandler,
		ReadTimeout:           0, // No read timeout for WebSocket
		WriteTimeout:          0, // No write timeout for WebSocket
		IdleTimeout:           0, // No idle timeout
		DisableKeepalive:      false,
		StreamRequestBody:     true,
		EnableIPValidation:    false,
		BodyLimit:             50 * 1024 * 1024, // 50MB for file uploads
	})
	
	cache.InitRedis()
	database.InitPostgres()

	// Recover from panics
	app.Use(recover.New())

	// Request logging
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	// Security headers
	app.Use(middleware.SecurityHeaders())

	// CORS configuration
	allowedOrigins := getAllowedOrigins()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/api/health", func(ctx *fiber.Ctx) error {
		server := os.Getenv("SERVER_NAME")
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "healthy", 
			"server": server,
		})
	})

	// Database health check (protected)
	app.Get("/api/health/database", auth.AuthorizationMiddleware, func(ctx *fiber.Ctx) error {
		healthStatus := database.GetHealthStatus()
		
		status := fiber.StatusOK
		if healthStatus["status"] == "unhealthy" {
			status = fiber.StatusServiceUnavailable
		} else if healthStatus["status"] == "degraded" {
			status = fiber.StatusPartialContent
		}
		
		return ctx.Status(status).JSON(healthStatus)
	})

	// Database metrics (protected)
	app.Get("/api/metrics/database", auth.AuthorizationMiddleware, func(ctx *fiber.Ctx) error {
		metrics := database.GetPerformanceMetrics()
		return ctx.Status(fiber.StatusOK).JSON(metrics)
	})

	// Auth endpoints with rate limiting
	app.Post("/api/auth/signup", middleware.AuthRateLimiter(), auth.SignUp)
	app.Post("/api/auth/login", middleware.AuthRateLimiter(), auth.Login)

	// Message history endpoint (protected with rate limiting)
	app.Get("/api/rooms/:roomId/messages", 
		middleware.GeneralRateLimiter(), 
		auth.AuthorizationMiddleware, 
		chat.GetRoomMessages)

	// File upload endpoint (protected with rate limiting)
	app.Post("/api/upload",
		middleware.GeneralRateLimiter(),
		auth.AuthorizationMiddleware,
		upload.UploadFile)

	// Get attachment metadata
	app.Get("/api/attachments/:id",
		middleware.GeneralRateLimiter(),
		auth.AuthorizationMiddleware,
		upload.GetAttachment)

	// Secure websocket connection
	app.Use("/ws", upgradeToWebSocket)
	app.Get("/ws/chat", websocket.New(chat.WebSocketHandler))

	port := ":" + config.Config.ServerPort
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}

// getAllowedOrigins returns allowed CORS origins based on environment
func getAllowedOrigins() []string {
	env := os.Getenv("NGINX_ENV")
	
	if env == "production" {
		// In production, specify exact origins
		origins := os.Getenv("ALLOWED_ORIGINS")
		if origins != "" {
			return strings.Split(origins, ",")
		}
		// Default production origins
		return []string{"https://yourdomain.com"}
	}
	
	// Development - allow localhost
	return []string{
		"http://localhost:8080",
		"http://localhost:3000",
		"http://localhost:3001",
	}
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}

// Authorize and Upgrate to websocket
func upgradeToWebSocket(context *fiber.Ctx) error {
	token := context.Query("token")
	if token == "" {
		log.Println("No token provided")
		return fiber.ErrUnauthorized
	}

	// Validate JWT token
	if err := auth.ValidateJWTToken(token); err != nil {
		log.Println("Error validating JWT token:", err)
		return fiber.ErrUnauthorized
	}

	userID, userName := auth.ParseJWTToken(token)
	if websocket.IsWebSocketUpgrade(context) {
		context.Locals("allowed", true)
		context.Locals("userID", userID)
		context.Locals("userName", userName)
		return context.Next()
	}
	return fiber.ErrUpgradeRequired
}
