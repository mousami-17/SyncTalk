package cache

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"realtime-chat/src/models"
)

var roomsMutex = &sync.Mutex{}
var subscribedRooms = make(map[string]bool)

// Store multiple callbacks per room to handle multiple connections to the same room
var roomCallbacks = make(map[string][]models.MessageHandlerCallbackType)
var roomCallbacksMutex = &sync.Mutex{}

// Track if listener is running
var listenerStarted = false
var listenerMutex = &sync.Mutex{}

func SubscribeToRoom(room string, callback models.MessageHandlerCallbackType) {
	// Start the listener ONCE when the package initializes
	listenerMutex.Lock()
	if !listenerStarted {
		listenerStarted = true
		listenerMutex.Unlock()
		log.Println("[PubSub] Starting message listener goroutine")
		go listenForMessages()
	} else {
		listenerMutex.Unlock()
	}

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	// Subscribe to room if not already subscribed
	if !subscribedRooms[room] {
		ctx := context.Background()
		err := PubSubConnection.Subscribe(ctx, room)
		if err != nil {
			log.Printf("[PubSub] Error subscribing to room %s: %v\n", room, err)
			return
		}
		subscribedRooms[room] = true
		log.Printf("[PubSub] Subscribed to room: %s\n", room)

		// Initialize the callbacks slice for this room
		roomCallbacksMutex.Lock()
		roomCallbacks[room] = []models.MessageHandlerCallbackType{callback}
		roomCallbacksMutex.Unlock()
	} else {
		// If already subscribed, append the new callback to the existing list
		roomCallbacksMutex.Lock()
		roomCallbacks[room] = append(roomCallbacks[room], callback)
		roomCallbacksMutex.Unlock()
		log.Printf("[PubSub] Added callback to existing room: %s\n", room)
	}
}

// Listen for messages on the single PubSub connection
func listenForMessages() {
	log.Println("[PubSub] 🎧 Message listener started and waiting for messages...")
	channel := PubSubConnection.Channel()
	for message := range channel {
		log.Printf("[PubSub] 📨 Received message from channel '%s': %s\n", message.Channel, message.Payload)
		var chatMessage models.Message
		err := json.Unmarshal([]byte(message.Payload), &chatMessage)
		if err != nil {
			log.Printf("[PubSub] Error decoding message from channel: %v\n", err)
			continue
		}

		// Get all callbacks for this specific room
		roomCallbacksMutex.Lock()
		callbacks, exists := roomCallbacks[message.Channel]
		callbackCount := len(callbacks)
		roomCallbacksMutex.Unlock()

		// Call all appropriate callback functions for the room
		if exists && callbackCount > 0 {
			log.Printf("[PubSub] 📢 Broadcasting to %d callback(s) for room '%s'\n", callbackCount, message.Channel)
			for i, callback := range callbacks {
				if callback != nil {
					log.Printf("[PubSub] Calling callback #%d for room '%s'\n", i+1, message.Channel)
					callback(message.Channel, &chatMessage)
				}
			}
		} else {
			log.Printf("[PubSub] ⚠️ No callbacks found for room '%s'\n", message.Channel)
		}
	}
	log.Println("[PubSub] ⚠️ Message listener stopped!")
}

func CheckAndUnsubscribeFromRoom(room string) {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	if subscribedRooms[room] {
		ctx := context.Background()
		key := "room:" + room
		members, _ := RedisClient.SCard(ctx, key).Result()

		// Unsubscribe from room if there are no members
		if members == 0 {
			err := PubSubConnection.Unsubscribe(ctx, room)
			if err != nil {
				log.Printf("Error unsubscribing from room %s: %v\n", room, err)
			}
			delete(subscribedRooms, room)
			
			// Remove callbacks for this room
			roomCallbacksMutex.Lock()
			delete(roomCallbacks, room)
			roomCallbacksMutex.Unlock()
		}
	}
}

func PublishMessage(room string, message *models.Message) {
	ctx := context.Background()
	msg, err := json.Marshal(message)
	if err != nil {
		log.Printf("[PubSub] Error marshalling message: %v\n", err)
		return
	}
	// Any connection from the pool can be used to publish messages
	err = RedisClient.Publish(ctx, room, msg).Err()
	if err != nil {
		log.Printf("[PubSub] Error publishing message to room %s: %v\n", room, err)
		return
	}
	log.Printf("[PubSub] Published message to channel '%s': %s\n", room, string(msg))
}

func GetAllRooms() []string {
	ctx := context.Background()
	keys, err := RedisClient.Keys(ctx, "room:*").Result()
	if err != nil {
		log.Println(err)
		return nil
	}
	return keys
}
