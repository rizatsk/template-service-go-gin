package main

import (
	"flag"
	"log"
	"service_go_gin/internal/config"
	"service_go_gin/pkg/database"
	"service_go_gin/pkg/database/seeder"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	action := flag.String("action", "up", "Action command: 'up', 'down', 'down-step', or 'seed-up'")
	flag.Parse()

	env := config.LoadEnvironment()

	switch *action {
	case "up":
		log.Println("Running migrate up...")
		database.RunMigrationsUp(env)
	case "down":
		log.Println("Running migrate down...")
		database.RunMigrationsDown(env)
	case "down-step":
		log.Println("Running migrate down...")
		database.RunMigrationsDownStep(env)
	case "seed-up":
		db_postgres := database.InitPostgres(env)
		seeder.RunAll(db_postgres, env)
	}
}
