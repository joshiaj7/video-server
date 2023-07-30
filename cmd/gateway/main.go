package main

import (
	"fmt"
	"log"
	"net/http"

	"video-server/internal/config"
)

func main() {
	cfg, err := config.NewGatewayServer()
	if err != nil {
		log.Fatalf("Load Gateway Server Failed: %v", err)
		return
	}

	conn, err := cfg.Database.DB()
	if err != nil {
		log.Fatalf("Get Database connection: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening to port 8080")
	http.ListenAndServe(":8080", cfg.Router)
}
