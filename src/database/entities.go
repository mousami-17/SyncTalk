package database

import (
	"time"

	"gorm.io/gorm"
)

type DBUser struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null;uniqueIndex"`
	Password string `gorm:"type:varchar(100);not null"`
}

type DBMessage struct {
	gorm.Model
	UserID    string `gorm:"type:varchar(100);not null;index"`
	RoomID    string `gorm:"type:varchar(100);not null;index:idx_room_timestamp"`
	Message   string `gorm:"type:varchar(1000);not null"`
	Timestamp *time.Time `gorm:"index:idx_room_timestamp"`
	
	// Media attachment fields
	FileURL   string `gorm:"type:text"`
	FileType  string `gorm:"type:varchar(50)"`
	FileName  string `gorm:"type:varchar(255)"`
	FileSize  int64  `gorm:"type:bigint"`
	Thumbnail string `gorm:"type:text"`
	Duration  int    `gorm:"type:int"`
	Width     int    `gorm:"type:int"`
	Height    int    `gorm:"type:int"`
}

type DBAttachment struct {
	gorm.Model
	MessageID  uint   `gorm:"index"`
	UserID     string `gorm:"type:varchar(100);not null;index"`
	RoomID     string `gorm:"type:varchar(100);not null;index"`
	FileURL    string `gorm:"type:varchar(500);not null"`
	FileType   string `gorm:"type:varchar(20);not null"` // "image", "video", "audio"
	FileName   string `gorm:"type:varchar(255)"`
	FileSize   int64
	Thumbnail  string `gorm:"type:varchar(500)"`
	Duration   int    // For audio/video in seconds
	Width      int    // For images/videos
	Height     int    // For images/videos
	MimeType   string `gorm:"type:varchar(100)"`
}
