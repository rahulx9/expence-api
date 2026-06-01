package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	store := NewStore("data.json")
	if err := store.Load(); err != nil {
		log.Fatalf("Failed to load store: %v", err)
	}

	r := gin.Default()
	registerRoutes(r, store)

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
