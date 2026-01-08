package main

import (
	"log"

	"github.com/Kredo15/task-board/services/api-gateway/internal/app"
)

func main() {
	application, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	if err := application.Run(); err != nil {
		log.Fatalf("Application runtime error: %v", err)
	}
}
