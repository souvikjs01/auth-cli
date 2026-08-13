package main

import (
	"fmt"
	"log"

	"github.com/souvikjs01/auth-cli/internals/config"
	"github.com/souvikjs01/auth-cli/internals/db"
	repository "github.com/souvikjs01/auth-cli/internals/repositories"
	"github.com/souvikjs01/auth-cli/internals/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	if err := db.Migrate(database); err != nil {
		log.Fatalf("migration error: %v", err)
	}

	userRepo := repository.NewUserRepository(database)
	_ = service.NewAuthService(
		userRepo,
		cfg.Auth.MaxLoginAttempts,
		cfg.Auth.LockoutDuration,
	)
	fmt.Println("Database connection successful")
	fmt.Println("Database migration successful")
}
