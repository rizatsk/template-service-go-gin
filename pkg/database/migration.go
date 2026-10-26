package database

import (
	"errors"
	"fmt"
	"log"
	"service_go_gin/internal/config"

	"github.com/golang-migrate/migrate/v4"
)

func MigrationInstance(env *config.Environment) *migrate.Migrate {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		env.DBPGUsername, env.DBPGPassword, env.DBPGHost, env.DBPGPort, env.DBPGName, env.DBPGSSLMode, env.DBPGSchema,
	)

	dbMigrate, err := migrate.New("file://database/migrations_postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed initial db migrate: %v", err)
	}
	return dbMigrate
}

func RunMigrationsUp(env *config.Environment) {
	dbMigrate := MigrationInstance(env)

	if err := dbMigrate.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed migrate up :%v", err)
	}
	log.Println("Success run migrate up")
}

func RunMigrationsDown(env *config.Environment) {
	dbMigrate := MigrationInstance(env)

	if err := dbMigrate.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed migrate down :%v", err)
	}
	log.Println("Success run migrate down")
}

func RunMigrationsDownStep(env *config.Environment) {
	dbMigrate := MigrationInstance(env)

	if err := dbMigrate.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed run mgirate down step :%v", err)
	}
	log.Println("Success run migrate down step")
}
