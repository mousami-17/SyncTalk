package chat

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"realtime-chat/src/database"
)

type MessageResponse struct {
	ID         uint   `json:"id"`
	Sender     string `json:"sender"`
	SenderName string `json:"senderName"`
	Room       string `json:"room"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
	Type       string `json:"type"`
	
	// Media attachment fields
	FileURL   string `json:"fileUrl,omitempty"`
	FileType  string `json:"fileType,omitempty"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Duration  int    `json:"duration,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}

type HistoryResponse struct {
	Messages   []MessageResponse `json:"messages"`
	TotalCount int64             `json:"totalCount"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
	HasMore    bool              `json:"hasMore"`
}

// GetRoomMessages fetches message history for a room with pagination
func GetRoomMessages(ctx *fiber.Ctx) error {
	roomID := ctx.Params("roomId")
	if roomID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Room ID is required",
		})
	}

	// Get pagination parameters
	limitStr := ctx.Query("limit", "50")
	offsetStr := ctx.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 50 // Default to 50, max 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Get total count of messages in room
	var totalCount int64
	if err := database.DB.Model(&database.DBMessage{}).
		Where("room_id = ?", roomID).
		Count(&totalCount).Error; err != nil {
		log.Printf("Error counting messages: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch message count",
		})
	}

	// Fetch messages with user information in a single query (prevents N+1)
	messagesWithUsers, err := database.GetMessagesWithUsers(roomID, limit, offset)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch messages",
		})
	}

	// Format messages for response
	messages := make([]MessageResponse, 0, len(messagesWithUsers))
	for _, msg := range messagesWithUsers {
		messages = append(messages, MessageResponse{
			ID:         msg.ID,
			Sender:     msg.Sender,
			SenderName: msg.SenderName,
			Room:       msg.Room,
			Content:    msg.Content,
			Timestamp:  msg.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			Type:       "chat_message",
			// Media fields
			FileURL:   msg.FileURL,
			FileType:  msg.FileType,
			FileName:  msg.FileName,
			FileSize:  msg.FileSize,
			Thumbnail: msg.Thumbnail,
			Duration:  msg.Duration,
			Width:     msg.Width,
			Height:    msg.Height,
		})
	}

	// Reverse the order so oldest messages come first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	hasMore := int64(offset+limit) < totalCount

	response := HistoryResponse{
		Messages:   messages,
		TotalCount: totalCount,
		Limit:      limit,
		Offset:     offset,
		HasMore:    hasMore,
	}

	return ctx.JSON(response)
}
