package database

import (
	"time"

	"gorm.io/gorm"
)

type DBMessage struct {
	gorm.Model
	UserID    string `gorm:"type:varchar(100);not null"`
	RoomID    string `gorm:"type:varchar(100);not null"`
	Message   string `gorm:"type:varchar(100);not null"`
	Timestamp *time.Time
	
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
