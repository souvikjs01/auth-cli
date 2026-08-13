package cmd

import (
	"log"

	"github.com/souvikjs01/auth-cli/internals/app"
	"github.com/souvikjs01/auth-cli/internals/config"
	"github.com/souvikjs01/auth-cli/internals/db"
	repository "github.com/souvikjs01/auth-cli/internals/repositories"
	"github.com/souvikjs01/auth-cli/internals/service"
)

var application *app.App

func initApp() {
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
	sessionRepo := repository.NewSessionRepository(database)

	authService := service.NewAuthService(
		userRepo,
		cfg.Auth.MaxLoginAttempts,
		cfg.Auth.LockoutDuration,
	)

	sessionService := service.NewSessionService(
		sessionRepo,
		cfg.Auth.SessionTimeout,
	)

	application = app.NewApp(
		authService,
		sessionService,
	)
}
