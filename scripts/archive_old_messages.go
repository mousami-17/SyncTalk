package main

import (
	"flag"
	"log"
	"time"

	"realtime-chat/src/config"
	"realtime-chat/src/database"
)

func main() {
	// Parse command line flags
	daysOld := flag.Int("days", 90, "Archive messages older than this many days")
	dryRun := flag.Bool("dry-run", false, "Perform a dry run without actually deleting")
	flag.Parse()

	log.Printf("Starting message archival process...")
	log.Printf("Archiving messages older than %d days", *daysOld)
	log.Printf("Dry run: %v", *dryRun)

	// Initialize database
	database.InitPostgres()

	// Calculate cutoff time
	cutoffDuration := time.Duration(*daysOld) * 24 * time.Hour

	if *dryRun {
		// Count messages that would be archived
		var count int64
		cutoffTime := time.Now().Add(-cutoffDuration)
		database.DB.Model(&database.DBMessage{}).
			Where("timestamp < ?", cutoffTime).
			Count(&count)
		
		log.Printf("Dry run: Would archive %d messages", count)
		return
	}

	// Archive old messages
	rowsAffected, err := database.ArchiveOldMessages(cutoffDuration)
	if err != nil {
		log.Fatalf("Failed to archive messages: %v", err)
	}

	log.Printf("Successfully archived %d messages", rowsAffected)
}
