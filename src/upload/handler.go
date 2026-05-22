package upload

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/cloudinary"
	"realtime-chat/src/database"
)

const (
	MaxFileSize = 50 * 1024 * 1024 // 50MB
)

// UploadFile handles file uploads
func UploadFile(ctx *fiber.Ctx) error {
	log.Println("========================================")
	log.Println("🚀 UPLOAD HANDLER CALLED - NEW CODE!")
	log.Println("========================================")
	
	// Get user ID from context (set by auth middleware)
	log.Printf("[Upload] Getting userID from context...")
	userIDRaw := ctx.Locals("userID")
	log.Printf("[Upload] userID raw value: %v (type: %T)", userIDRaw, userIDRaw)
	
	if userIDRaw == nil {
		log.Println("[Upload] ERROR: userID is nil in context")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not authenticated",
		})
	}
	
	userID, ok := userIDRaw.(string)
	if !ok {
		log.Printf("[Upload] ERROR: userID type assertion failed, got type: %T", userIDRaw)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user context",
		})
	}
	
	log.Printf("[Upload] Successfully got userID: %s", userID)
	
	// Get room ID from form
	roomID := ctx.FormValue("roomId")
	if roomID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Room ID is required",
		})
	}

	// Get file from form
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Printf("Error getting file: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file provided",
		})
	}

	// Check file size
	if file.Size > MaxFileSize {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File size exceeds 50MB limit",
		})
	}

	// Get MIME type
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type",
		})
	}

	// Validate file type
	fileType := cloudinary.GetFileType(mimeType)
	if fileType != "image" && fileType != "video" && fileType != "audio" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Only images, videos, and audio files are allowed",
		})
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to process file",
		})
	}
	defer fileContent.Close()

	// Upload to Cloudinary
	resourceType := cloudinary.GetResourceType(mimeType)
	cloudinaryResp, err := cloudinary.UploadFile(fileContent, file.Filename, resourceType)
	if err != nil {
		log.Printf("Error uploading to Cloudinary: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to upload file",
		})
	}

	// Generate thumbnail
	thumbnail := cloudinary.GenerateThumbnailURL(cloudinaryResp.PublicID, resourceType)

	// Save to database
	attachment := database.DBAttachment{
		UserID:    userID,
		RoomID:    roomID,
		FileURL:   cloudinaryResp.SecureURL,
		FileType:  fileType,
		FileName:  file.Filename,
		FileSize:  file.Size,
		Thumbnail: thumbnail,
		Duration:  int(cloudinaryResp.Duration),
		Width:     cloudinaryResp.Width,
		Height:    cloudinaryResp.Height,
		MimeType:  mimeType,
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		log.Printf("Error saving attachment to database: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save file metadata",
		})
	}

	log.Printf("File uploaded successfully: %s (ID: %d)", file.Filename, attachment.ID)

	// Return file info
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":     false,
		"message":   "File uploaded successfully",
		"id":        attachment.ID,
		"fileUrl":   cloudinaryResp.SecureURL,
		"fileType":  fileType,
		"fileName":  file.Filename,
		"fileSize":  file.Size,
		"thumbnail": thumbnail,
		"duration":  int(cloudinaryResp.Duration),
		"width":     cloudinaryResp.Width,
		"height":    cloudinaryResp.Height,
	})
}

// GetAttachment retrieves attachment metadata
func GetAttachment(ctx *fiber.Ctx) error {
	attachmentID := ctx.Params("id")
	if attachmentID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Attachment ID is required",
		})
	}

	var attachment database.DBAttachment
	if err := database.DB.First(&attachment, attachmentID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Attachment not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"error":     false,
		"id":        attachment.ID,
		"fileUrl":   attachment.FileURL,
		"fileType":  attachment.FileType,
		"fileName":  attachment.FileName,
		"fileSize":  attachment.FileSize,
		"thumbnail": attachment.Thumbnail,
		"duration":  attachment.Duration,
		"width":     attachment.Width,
		"height":    attachment.Height,
	})
}
