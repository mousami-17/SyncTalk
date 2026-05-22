package chat

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"realtime-chat/src/cache"
	"realtime-chat/src/models"
)

// User presence status
const (
	StatusOnline  = "online"
	StatusOffline = "offline"
	StatusTyping  = "typing"
)

// Presence TTL in Redis (5 minutes)
const PresenceTTL = 5 * time.Minute

// Typing indicator TTL (5 seconds)
const TypingTTL = 5 * time.Second

// SetUserOnline marks a user as online in a room
func SetUserOnline(userID, userName, roomID string) error {
	ctx := context.Background()
	key := "presence:room:" + roomID
	
	presenceData := map[string]interface{}{
		"userId":    userID,
		"userName":  userName,
		"status":    StatusOnline,
		"timestamp": time.Now().Unix(),
	}

	data, err := json.Marshal(presenceData)
	if err != nil {
		return err
	}

	// Add user to room presence set with TTL
	if err := cache.RedisClient.HSet(ctx, key, userID, data).Err(); err != nil {
		return err
	}

	// Set expiration on the hash key
	cache.RedisClient.Expire(ctx, key, PresenceTTL)

	// Broadcast presence update
	broadcastPresenceUpdate(roomID, userID, userName, StatusOnline)
	
	return nil
}

// SetUserOffline marks a user as offline in a room
func SetUserOffline(userID, roomID string) error {
	ctx := context.Background()
	key := "presence:room:" + roomID

	// Get user data before removing
	userData, err := cache.RedisClient.HGet(ctx, key, userID).Result()
	if err == nil {
		var presenceData map[string]interface{}
		if err := json.Unmarshal([]byte(userData), &presenceData); err == nil {
			userName := presenceData["userName"].(string)
			broadcastPresenceUpdate(roomID, userID, userName, StatusOffline)
		}
	}

	// Remove user from presence set
	return cache.RedisClient.HDel(ctx, key, userID).Err()
}

// GetRoomUsers returns all online users in a room
func GetRoomUsers(roomID string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	key := "presence:room:" + roomID

	users := []map[string]interface{}{}
	
	// Get all users in room
	usersData, err := cache.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return users, err
	}

	for _, userData := range usersData {
		var presenceData map[string]interface{}
		if err := json.Unmarshal([]byte(userData), &presenceData); err == nil {
			users = append(users, presenceData)
		}
	}

	return users, nil
}

// SetUserTyping marks a user as typing in a room
func SetUserTyping(userID, userName, roomID string) error {
	ctx := context.Background()
	key := "typing:room:" + roomID

	typingData := map[string]interface{}{
		"userId":   userID,
		"userName": userName,
	}

	data, err := json.Marshal(typingData)
	if err != nil {
		return err
	}

	// Add user to typing set with short TTL
	if err := cache.RedisClient.HSet(ctx, key, userID, data).Err(); err != nil {
		return err
	}

	// Set short expiration
	cache.RedisClient.Expire(ctx, key, TypingTTL)

	// Broadcast typing indicator
	broadcastTypingIndicator(roomID, userID, userName, true)

	return nil
}

// RemoveUserTyping removes typing indicator for a user
func RemoveUserTyping(userID, roomID string) error {
	ctx := context.Background()
	key := "typing:room:" + roomID

	// Get user data before removing
	userData, err := cache.RedisClient.HGet(ctx, key, userID).Result()
	if err == nil {
		var typingData map[string]interface{}
		if err := json.Unmarshal([]byte(userData), &typingData); err == nil {
			userName := typingData["userName"].(string)
			broadcastTypingIndicator(roomID, userID, userName, false)
		}
	}

	return cache.RedisClient.HDel(ctx, key, userID).Err()
}

// GetTypingUsers returns all users currently typing in a room
func GetTypingUsers(roomID string) ([]map[string]interface{}, error) {
	ctx := context.Background()
	key := "typing:room:" + roomID

	users := []map[string]interface{}{}
	
	usersData, err := cache.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return users, err
	}

	for _, userData := range usersData {
		var typingData map[string]interface{}
		if err := json.Unmarshal([]byte(userData), &typingData); err == nil {
			users = append(users, typingData)
		}
	}

	return users, nil
}

// UpdateUserPresence updates user presence with heartbeat
func UpdateUserPresence(userID, userName, roomID string) error {
	return SetUserOnline(userID, userName, roomID)
}

// broadcastPresenceUpdate sends presence update to all users in room
func broadcastPresenceUpdate(roomID, userID, userName, status string) {
	// Get all current users in the room
	roomUsers, err := GetRoomUsers(roomID)
	if err != nil {
		log.Printf("Failed to get room users for presence update: %v", err)
		roomUsers = []map[string]interface{}{}
	}

	// Convert users to the format expected by frontend
	users := make([]map[string]interface{}, 0, len(roomUsers))
	for _, user := range roomUsers {
		users = append(users, map[string]interface{}{
			"id":     user["userId"],
			"name":   user["userName"],
			"status": user["status"],
		})
	}

	// Create a custom message with users array
	messageWithUsers := map[string]interface{}{
		"type":       "presence_update",
		"room":       roomID,
		"sender":     userID,
		"senderName": userName,
		"content":    status,
		"users":      users,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	// Publish the enhanced message directly to Redis
	ctx := context.Background()
	if data, err := json.Marshal(messageWithUsers); err == nil {
		cache.RedisClient.Publish(ctx, roomID, string(data))
		log.Printf("Presence update: %s is %s in room %s (total users: %d)", userName, status, roomID, len(users))
	} else {
		log.Printf("Failed to marshal presence update: %v", err)
	}
}

// broadcastTypingIndicator sends typing indicator to all users in room
func broadcastTypingIndicator(roomID, userID, userName string, isTyping bool) {
	status := "stopped_typing"
	if isTyping {
		status = "typing"
	}

	message := models.Message{
		Type:       "typing_indicator",
		Room:       roomID,
		Sender:     userID,
		SenderName: userName,
		Content:    status,
	}

	cache.PublishMessage(roomID, &message)
}
