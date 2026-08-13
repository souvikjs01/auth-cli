package main

import (
	"fmt"
	"log"

	"github.com/souvikjs01/auth-cli/internals/config"
	"github.com/souvikjs01/auth-cli/internals/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	_, err = db.Connect(cfg)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	fmt.Println("Database connection successful")
}
