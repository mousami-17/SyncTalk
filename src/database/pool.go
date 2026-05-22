package database

import (
	"log"
	"time"
)

// ConnectionPoolConfig holds database connection pool settings
type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// GetOptimalPoolConfig returns optimal connection pool settings based on environment
func GetOptimalPoolConfig() ConnectionPoolConfig {
	// Default configuration for production
	config := ConnectionPoolConfig{
		MaxOpenConns:    25,              // Maximum number of open connections
		MaxIdleConns:    10,              // Maximum number of idle connections
		ConnMaxLifetime: 1 * time.Hour,   // Maximum lifetime of a connection
		ConnMaxIdleTime: 10 * time.Minute, // Maximum idle time before closing
	}

	// Adjust based on expected load
	// For high-traffic applications, increase these values
	// For low-traffic applications, decrease to save resources

	return config
}

// ConfigureConnectionPool sets up the connection pool with optimal settings
func ConfigureConnectionPool() {
	config := GetOptimalPoolConfig()

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	log.Printf("Set MaxOpenConns to %d", config.MaxOpenConns)

	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	log.Printf("Set MaxIdleConns to %d", config.MaxIdleConns)

	// Set maximum lifetime of a connection
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	log.Printf("Set ConnMaxLifetime to %v", config.ConnMaxLifetime)

	// Set maximum idle time of a connection
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	log.Printf("Set ConnMaxIdleTime to %v", config.ConnMaxIdleTime)

	log.Println("Database connection pool configured successfully")
}

// GetPoolStats returns current connection pool statistics
func GetPoolStats() map[string]interface{} {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Failed to get database instance: %v", err)
		return nil
	}

	stats := sqlDB.Stats()

	return map[string]interface{}{
		"maxOpenConnections":  stats.MaxOpenConnections,
		"openConnections":     stats.OpenConnections,
		"inUse":               stats.InUse,
		"idle":                stats.Idle,
		"waitCount":           stats.WaitCount,
		"waitDuration":        stats.WaitDuration.String(),
		"maxIdleClosed":       stats.MaxIdleClosed,
		"maxLifetimeClosed":   stats.MaxLifetimeClosed,
		"maxIdleTimeClosed":   stats.MaxIdleTimeClosed,
	}
}
