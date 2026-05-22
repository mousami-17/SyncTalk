package chat

import (
	"sync"

	"github.com/gofiber/contrib/websocket"
)

var connectionMutex = &sync.Mutex{}
var UserConnections = make(map[string]*websocket.Conn)

// Mutex for each WebSocket connection to prevent concurrent writes
var writeMutexes = make(map[string]*sync.Mutex)
var writeMutexesLock = &sync.Mutex{}

// Register a new connection in the Map
func AddConnection(userID string, conn *websocket.Conn) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	UserConnections[userID] = conn
	
	// Create a write mutex for this connection
	writeMutexesLock.Lock()
	writeMutexes[userID] = &sync.Mutex{}
	writeMutexesLock.Unlock()
}

// Remove a connection from the Map
func RemoveConnection(userID string) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	delete(UserConnections, userID)
	
	// Remove the write mutex
	writeMutexesLock.Lock()
	delete(writeMutexes, userID)
	writeMutexesLock.Unlock()
}

// Get a connection from the Map
func GetConnection(userID string) (*websocket.Conn, bool) {
	connectionMutex.Lock()
	defer connectionMutex.Unlock()
	conn, ok := UserConnections[userID]
	return conn, ok
}

// SafeWriteJSON writes JSON to a WebSocket connection with mutex protection
func SafeWriteJSON(conn *websocket.Conn, v interface{}) error {
	// Find the userID for this connection
	connectionMutex.Lock()
	var userID string
	for id, c := range UserConnections {
		if c == conn {
			userID = id
			break
		}
	}
	connectionMutex.Unlock()
	
	if userID == "" {
		return conn.WriteJSON(v)
	}
	
	// Get the write mutex for this connection
	writeMutexesLock.Lock()
	mutex, exists := writeMutexes[userID]
	writeMutexesLock.Unlock()
	
	if !exists {
		return conn.WriteJSON(v)
	}
	
	// Lock before writing
	mutex.Lock()
	defer mutex.Unlock()
	return conn.WriteJSON(v)
}
