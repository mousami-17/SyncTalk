package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"realtime-chat/src/cache"
	"realtime-chat/src/kafka"
	"realtime-chat/src/models"
	"realtime-chat/src/utils"
)

var roomMutex = &sync.Mutex{}
var subscribedRooms = make(map[string]bool)

func JoinRoom(room string, user *models.User) {
	key := "room:" + room
	if err := addUserToRoomInRedis(key, user); err != nil {
		log.Println("Failed to add user to room:", err)
		utils.SendErrorMessage(user.Connection, "Unable to join room")
		return
	}

	// Only subscribe once per room (not per user)
	roomMutex.Lock()
	if !subscribedRooms[room] {
		subscribedRooms[room] = true
		roomMutex.Unlock()
		
		cache.SubscribeToRoom(room, func(room string, message *models.Message) {
			BroadcastToRoom(room, *message)
		})
	} else {
		roomMutex.Unlock()
	}

	// Set user as online in the room
	if err := SetUserOnline(user.ID, user.Name, room); err != nil {
		log.Println("Failed to set user online:", err)
	}

	message := models.Message{
		Sender:     user.ID,
		SenderName: user.Name,
		Room:       room,
		Type:       "join_room",
		Server:     os.Getenv("SERVER_NAME"),
	}
	cache.PublishMessage(room, &message)

	// Get current room users and send to the joining user
	roomUsers, err := GetRoomUsers(room)
	if err != nil {
		log.Println("Failed to get room users:", err)
	}

	log.Printf("User %s joined room %s\n", user.ID, room)
	response := map[string]interface{}{
		"type":    "join_room",
		"room":    room,
		"success": true,
		"users":   roomUsers,
	}
	if err := SafeWriteJSON(user.Connection, response); err != nil {
		log.Println("Error sending join room response:", err)
	}
}

func BroadcastToRoom(room string, message models.Message) {
	log.Printf("[Broadcast] 🚀 BroadcastToRoom called for room '%s', message type: %s", room, message.Type)
	key := "room:" + room
	members := getAllMembersInRoom(key)
	log.Printf("[Broadcast] Broadcasting to room '%s' with %d members: %v", room, len(members), members)
	
	for _, userID := range members {
		// Get the websocket connection for the user from the local map
		if conn, exists := GetConnection(userID); exists {
			log.Printf("[Broadcast] Sending message to user %s (connection exists)", userID)
			// Send the message to the user with mutex protection
			if err := SafeWriteJSON(conn, message); err != nil {
				log.Printf("[Broadcast] Error sending message to user %s: %v\n", userID, err)
				conn.Close()
				RemoveConnection(userID)
			} else {
				log.Printf("[Broadcast] Successfully sent message to user %s", userID)
			}
		} else {
			log.Printf("[Broadcast] ⚠️ No connection found for user %s", userID)
		}
	}
}

func SendMessageToRoom(message models.Message, user *models.User) {
	// Ensure we're subscribed to the room before sending
	roomMutex.Lock()
	if !subscribedRooms[message.Room] {
		subscribedRooms[message.Room] = true
		roomMutex.Unlock()
		
		cache.SubscribeToRoom(message.Room, func(room string, message *models.Message) {
			BroadcastToRoom(room, *message)
		})
	} else {
		roomMutex.Unlock()
	}
	
	cache.PublishMessage(message.Room, &message)
	kafka.PublishMessage(message)

	if strings.HasPrefix(message.Content, "@superchat") {
		response := GetResponseFromLLM(message.Content)
		log.Printf("Response from LLM")
		llmMessage := models.Message{
			Content: response,
			Sender: "SuperChat",
			SenderName: "SuperChat",
			Room: message.Room,
			Type: message.Type,
			Server: os.Getenv("SERVER_NAME"),
		}
		// Publish to both Redis (for real-time) and Kafka (for persistence)
		cache.PublishMessage(message.Room, &llmMessage)
		kafka.PublishMessage(llmMessage)
	}
}


func GetResponseFromLLM(prompt string) string {
	// Get Gemini API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Printf("GEMINI_API_KEY not set")
		return "AI service not configured. Please set GEMINI_API_KEY."
	}

	// Remove @superchat prefix from prompt
	cleanPrompt := strings.TrimPrefix(prompt, "@superchat")
	cleanPrompt = strings.TrimSpace(cleanPrompt)

	// Gemini API endpoint - using gemini-1.5-flash for better free tier limits
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-2.5-flash:generateContent?key=%s", apiKey)


	// Build request body for Gemini APIs
	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{
						"text": cleanPrompt,
					},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": 150,
			"topP":            0.95,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return "Error preparing AI request"
	}

	// Make HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		return "Error connecting to AI service"
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
		
		// Check for rate limit errors
		if resp.StatusCode == 429 || strings.Contains(string(body), "quota") || strings.Contains(string(body), "RESOURCE_EXHAUSTED") {
			return "⚠️ AI service rate limit reached. Please try again in a minute."
		}
		
		return "AI service returned an error. Please try again later."
	}

	// Parse response
	var geminiResponse struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResponse); err != nil {
		log.Printf("Error decoding Gemini response: %v", err)
		return "Error processing AI response"
	}

	// Extract text from response
	if len(geminiResponse.Candidates) > 0 && 
	   len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		responseText := geminiResponse.Candidates[0].Content.Parts[0].Text
		log.Printf("✨ Gemini response: %s", responseText)
		return responseText
	}

	log.Printf("No response from Gemini")
	return "AI did not generate a response"
}


func LeaveRoom(room string, user *models.User) {
	key := "room:" + room
	if err := removeUserFromRoomInRedis(key, user); err != nil {
		log.Println("Failed to remove user from room:", err)
		utils.SendErrorMessage(user.Connection, "Unable to leave room")
		return
	}

	// Set user as offline in the room
	if err := SetUserOffline(user.ID, room); err != nil {
		log.Println("Failed to set user offline:", err)
	}

	// Remove typing indicator if any
	RemoveUserTyping(user.ID, room)

	cache.CheckAndUnsubscribeFromRoom(room)
	response := map[string]interface{}{
		"type":    "leave_room",
		"room":    room,
		"success": true,
	}
	if err := SafeWriteJSON(user.Connection, response); err != nil {
		log.Println("Error sending leave room response:", err)
	}
}

func LeaveAllRooms(user *models.User) {
	for _, room := range cache.GetAllRooms() {
		isMember := isUserInRoom(room, user)
		if isMember {
			removeUserFromRoomInRedis(room, user)
			// Set user as offline
			SetUserOffline(user.ID, room)
			// Remove typing indicator
			RemoveUserTyping(user.ID, room)
		}
	}
}

// Helper methods
func addUserToRoomInRedis(room string, user *models.User) error {
	ctx := context.Background()
	_, err := cache.RedisClient.SAdd(ctx, room, user.ID).Result()
	return err
}

func removeUserFromRoomInRedis(room string, user *models.User) error {
	ctx := context.Background()
	_, err := cache.RedisClient.SRem(ctx, room, user.ID).Result()
	return err
}

func isUserInRoom(room string, user *models.User) bool {
	ctx := context.Background()
	isMember, err := cache.RedisClient.SIsMember(ctx, room, user.ID).Result()
	if err != nil {
		log.Println(err)
		return false
	}
	return isMember
}

func getAllMembersInRoom(room string) []string {
	ctx := context.Background()
	members, err := cache.RedisClient.SMembers(ctx, room).Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return members
}
