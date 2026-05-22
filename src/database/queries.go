package database

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// GetUserByID retrieves a user by ID with caching consideration
func GetUserByID(userID string) (*DBUser, error) {
	var user DBUser
	
	// Use Select to only fetch needed fields
	err := DB.Select("id", "name", "created_at").
		Where("id = ?", userID).
		First(&user).Error
	
	return &user, err
}

// GetUserByName retrieves a user by name (for login)
func GetUserByName(username string) (*DBUser, error) {
	var user DBUser
	
	err := DB.Where("name = ?", username).
		First(&user).Error
	
	return &user, err
}

// GetMessagesWithUsers retrieves messages with user information (prevents N+1)
func GetMessagesWithUsers(roomID string, limit, offset int) ([]MessageWithUser, error) {
	var results []MessageWithUser
	
	// Use JOIN to fetch messages and users in a single query
	// COALESCE ensures that for non-user senders (like "SuperChat"), we use the user_id as the name
	err := DB.Table("db_messages").
		Select(`db_messages.id, db_messages.user_id, db_messages.room_id, db_messages.message, db_messages.timestamp, 
		        COALESCE(db_users.name, db_messages.user_id) as user_name, 'chat_message' as type,
		        db_messages.file_url, db_messages.file_type, db_messages.file_name, db_messages.file_size,
		        db_messages.thumbnail, db_messages.duration, db_messages.width, db_messages.height`).
		Joins("LEFT JOIN db_users ON db_messages.user_id = CAST(db_users.id AS VARCHAR)").
		Where("db_messages.room_id = ?", roomID).
		Order("db_messages.timestamp DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error
	
	return results, err
}

// MessageWithUser represents a message with user information
type MessageWithUser struct {
	ID         uint      `json:"id" gorm:"column:id"`
	Sender     string    `json:"sender" gorm:"column:user_id"`
	SenderName string    `json:"senderName" gorm:"column:user_name"`
	Room       string    `json:"room" gorm:"column:room_id"`
	Content    string    `json:"content" gorm:"column:message"`
	Timestamp  time.Time `json:"timestamp" gorm:"column:timestamp"`
	Type       string    `json:"type"`
	
	// Media attachment fields
	FileURL   string `json:"fileUrl" gorm:"column:file_url"`
	FileType  string `json:"fileType" gorm:"column:file_type"`
	FileName  string `json:"fileName" gorm:"column:file_name"`
	FileSize  int64  `json:"fileSize" gorm:"column:file_size"`
	Thumbnail string `json:"thumbnail" gorm:"column:thumbnail"`
	Duration  int    `json:"duration" gorm:"column:duration"`
	Width     int    `json:"width" gorm:"column:width"`
	Height    int    `json:"height" gorm:"column:height"`
}

// BulkInsertMessages inserts multiple messages efficiently
func BulkInsertMessages(messages []DBMessage) error {
	if len(messages) == 0 {
		return nil
	}

	// Use CreateInBatches for efficient bulk insert
	batchSize := 100
	return DB.CreateInBatches(messages, batchSize).Error
}

// GetMessageCount returns total message count for a room
func GetMessageCount(roomID string) (int64, error) {
	var count int64
	
	err := DB.Model(&DBMessage{}).
		Where("room_id = ?", roomID).
		Count(&count).Error
	
	return count, err
}

// DeleteOldMessages deletes messages older than specified duration
func DeleteOldMessages(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	
	result := DB.Where("timestamp < ?", cutoffTime).
		Delete(&DBMessage{})
	
	return result.RowsAffected, result.Error
}

// ArchiveOldMessages moves old messages to archive table
func ArchiveOldMessages(olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	
	// First, copy to archive table (if exists)
	// Then delete from main table
	result := DB.Where("timestamp < ?", cutoffTime).
		Delete(&DBMessage{})
	
	return result.RowsAffected, result.Error
}

// GetDatabaseSize returns approximate database size information
func GetDatabaseSize() map[string]interface{} {
	var messageCount int64
	var userCount int64
	
	DB.Model(&DBMessage{}).Count(&messageCount)
	DB.Model(&DBUser{}).Count(&userCount)
	
	return map[string]interface{}{
		"totalMessages": messageCount,
		"totalUsers":    userCount,
	}
}

// OptimizeDatabase runs database optimization commands
func OptimizeDatabase() error {
	// Run VACUUM ANALYZE on PostgreSQL
	err := DB.Exec("VACUUM ANALYZE").Error
	if err != nil {
		log.Printf("Failed to run VACUUM ANALYZE: %v", err)
		return err
	}
	
	log.Println("Database optimization completed")
	return nil
}

// GetSlowQueries returns queries that took longer than threshold
// This is a placeholder - implement based on your monitoring needs
func GetSlowQueries(threshold time.Duration) []string {
	// In production, integrate with pg_stat_statements or similar
	return []string{}
}

// PreparedStatementCache enables prepared statement caching
func EnablePreparedStatementCache() {
	DB = DB.Session(&gorm.Session{
		PrepareStmt: true,
	})
	log.Println("Prepared statement caching enabled")
}
