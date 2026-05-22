package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"realtime-chat/src/config"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var DB *gorm.DB
var postgresConfig PostgresConfig

func init() {
	postgresConfig = PostgresConfig{
		Host:     config.Config.PostgresHost,
		Port:     config.Config.PostgresPort,
		User:     config.Config.PostgresUser,
		Password: config.Config.PostgresPassword,
		Database: config.Config.PostgresDatabase,
	}
}

func InitPostgres() {
	var err error
	dsn := "host=" + postgresConfig.Host + " user=" + postgresConfig.User + " password=" + postgresConfig.Password + " dbname=" + postgresConfig.Database + " port=" + postgresConfig.Port + " sslmode=disable"
	
	// Configure GORM with optimizations
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // Enable prepared statement caching
		// SkipDefaultTransaction: true, // Uncomment for better performance if you don't need transactions
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Configure connection pool
	ConfigureConnectionPool()

	// Migrate the schema
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err = DB.AutoMigrate(&DBUser{}, &DBMessage{}, &DBAttachment{})
	if err != nil {
		// Log the error but don't fail - columns might already exist
		log.Printf("Database migration warning: %v", err)
		log.Println("Continuing with existing schema...")
	} else {
		log.Println("Database migration completed successfully")
	}

	// Log initial pool stats
	stats := GetPoolStats()
	log.Printf("Initial pool stats: %+v", stats)
}
