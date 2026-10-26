package database

import (
	"fmt"
	"log"
	"service_go_gin/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(environemnt *config.Environment) *gorm.DB {
	// DSN (data source name)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		environemnt.DBPGHost, environemnt.DBPGUsername, environemnt.DBPGPassword, environemnt.DBPGName, environemnt.DBPGPort, environemnt.DBPGSSLMode,
	)

	db_postgres, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot connect to database Postgres: %v", err)
	}

	// Connceetion pooling
	sqlDB, err := db_postgres.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan sql.DB instance: %v", err)
	}

	// Limit max conenction open to DB
	sqlDB.SetMaxOpenConns(15)

	// Limit max connection idle ready for reused
	sqlDB.SetMaxIdleConns(5)

	// Max duration connection idle before closed
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	// Max duration age connection (prevent memory issue leak or connection expired)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	log.Println("Success connection db Postgres with GORM, and pooling connection")
	return db_postgres
}
