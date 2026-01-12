// Package main implements a RESTful API server for managing a collection of albums.
// The server uses the Gin web framework and stores data in memory.
package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

const (
	serverPort = "localhost:8080"
)

// main initializes the Gin router, registers all API routes, and starts the HTTP server.
// The server listens on localhost:8080 and provides RESTful endpoints for album management.
func main() {
	router := gin.Default()

	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.PATCH("/albums/:id", patchAlbumByID)
	router.GET("/", healthCheck)

	log.Println("Starting Album API server...")
	log.Printf("Server listening on http://%s", serverPort)
	log.Println("Available endpoints:")
	log.Println("  GET    /albums      - List all albums")
	log.Println("  GET    /albums/:id  - Get album by ID")
	log.Println("  POST   /albums      - Create new album")
	log.Println("  DELETE /albums/:id  - Delete album by ID")
	log.Println("  PATCH  /albums/:id  - Update album by ID")
	log.Println("  GET    /            - Health check")

	if err := router.Run(serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
