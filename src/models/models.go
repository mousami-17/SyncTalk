package models

import (
	"github.com/gofiber/contrib/websocket"
)

type User struct {
	ID         string
	Name       string
	Connection *websocket.Conn
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MessageHandlerCallbackType func(room string, message *Message)

type Message struct {
	Sender     string `json:"sender"`
	SenderName string `json:"senderName"`
	Room       string `json:"room"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Server     string `json:"server"`
	Timestamp  string `json:"timestamp,omitempty"`
	ID         uint   `json:"id,omitempty"`
	
	// Media fields (optional)
	FileURL    string `json:"fileUrl,omitempty"`
	FileType   string `json:"fileType,omitempty"`    // "image", "video", "audio"
	FileName   string `json:"fileName,omitempty"`
	FileSize   int64  `json:"fileSize,omitempty"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Duration   int    `json:"duration,omitempty"`    // For audio/video in seconds
	Width      int    `json:"width,omitempty"`       // For images/videos
	Height     int    `json:"height,omitempty"`      // For images/videos
}

type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
