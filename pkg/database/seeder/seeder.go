package seeder

import (
	"service_go_gin/internal/config"

	"gorm.io/gorm"
)

func RunAll(db_postgres *gorm.DB, environment *config.Environment) {
	SeedUsers(db_postgres, environment)
}
