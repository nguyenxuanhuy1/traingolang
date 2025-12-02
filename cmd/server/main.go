package main

import (
	"log"

	"traingolang/internal/api/router"
	"traingolang/internal/config"
)

func main() {
	config.ConnectDB()
	r := router.SetupRouter()

	log.Println("Server is running on :8080")
	r.Run(":8080")
}
