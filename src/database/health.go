package database

import (
	"context"
	"log"
	"time"
)

// HealthCheck performs a database health check
func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// GetHealthStatus returns detailed database health information
func GetHealthStatus() map[string]interface{} {
	health := map[string]interface{}{
		"status": "healthy",
	}

	// Check connection
	if err := HealthCheck(); err != nil {
		health["status"] = "unhealthy"
		health["error"] = err.Error()
		return health
	}

	// Get pool stats
	poolStats := GetPoolStats()
	health["poolStats"] = poolStats

	// Get database size
	dbSize := GetDatabaseSize()
	health["databaseSize"] = dbSize

	// Check if we can execute a simple query
	var result int
	if err := DB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		health["status"] = "degraded"
		health["queryError"] = err.Error()
	}

	return health
}

// MonitorConnectionPool logs connection pool statistics periodically
func MonitorConnectionPool(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		stats := GetPoolStats()
		log.Printf("Connection Pool Stats: %+v", stats)

		// Alert if pool is exhausted
		if stats["inUse"].(int) >= stats["maxOpenConnections"].(int) {
			log.Printf("WARNING: Connection pool exhausted! Consider increasing MaxOpenConns")
		}

		// Alert if wait count is high
		if stats["waitCount"].(int64) > 100 {
			log.Printf("WARNING: High wait count (%d). Connections are being waited for frequently", stats["waitCount"])
		}
	}
}

// GetPerformanceMetrics returns database performance metrics
func GetPerformanceMetrics() map[string]interface{} {
	metrics := make(map[string]interface{})

	// Measure query performance
	start := time.Now()
	var count int64
	DB.Model(&DBMessage{}).Count(&count)
	queryDuration := time.Since(start)

	metrics["sampleQueryDuration"] = queryDuration.String()
	metrics["messageCount"] = count

	// Get pool stats
	metrics["poolStats"] = GetPoolStats()

	return metrics
}
