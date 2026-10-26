package seeder

import (
	"fmt"
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/internal/domain"
	"uuid"

	"gorm.io/gorm"
)

func SeedUsers(db_postgres *gorm.DB, environment *config.Environment) {
	log.Println("Check first data (seeder) users...")

	adminEmail := "rjshubz@gmail.com"
	admin := domain.User{
		ID:    uuid.NewV7().String(),
		Name:  "Rizat sakmir",
		Email: adminEmail,
		Role:  "admin",
	}

	if err := admin.HashPassword(environment.PWDefaultAdmin); err != nil {
		log.Fatalf("Failed hashed passowrd: %v", err)
	}

	var user domain.User
	result := db_postgres.Where("email=?", adminEmail).Attrs(admin).FirstOrCreate(&user)
	if result.Error != nil {
		log.Fatalf("Failed execute seeding user admin: %v", result.Error)
		return
	}

	if result.RowsAffected == 1 {
		message_log := fmt.Sprintf("Success create user admin: %s", adminEmail)
		log.Println(message_log)
	} else {
		log.Println("User have existing")
	}
}
