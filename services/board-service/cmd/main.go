package main

import (
	"log"

	"github.com/Kredo15/task-board/services/board-service/internal/app"
)

func main() {

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
