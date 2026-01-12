package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getAlbums handles GET /albums requests.
// Returns all albums in the collection as a JSON array with HTTP 200 status.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// healthCheck handles GET / requests.
// Returns the server health status as JSON with HTTP 200 status.
// Used for monitoring and load balancer health checks.
func healthCheck(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "album-api",
		"version": "1.0.0",
	})
}

// postAlbums handles POST /albums requests.
// Creates a new album with an auto-generated UUID. Validates all required fields.
// Returns the created album as JSON with HTTP 201 status on success,
// or HTTP 400 with error details if validation fails.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	if err := c.ShouldBindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON",
			"details": err.Error(),
		})
		return
	}

	if errMsg := validateTitle(newAlbum.Title, true); errMsg != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	if errMsg := validateArtist(newAlbum.Artist, true); errMsg != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}
	if errMsg := validatePrice(newAlbum.Price, true); errMsg != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	newAlbum.ID = uuid.New().String()
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID handles GET /albums/:id requests.
// Returns the album with the specified ID as JSON with HTTP 200 status.
// Returns HTTP 404 if the album is not found.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// deleteAlbumByID handles DELETE /albums/:id requests.
// Deletes the album with the specified ID and returns the deleted album as JSON with HTTP 200 status.
// Returns HTTP 404 if the album is not found.
func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// patchAlbumByID handles PATCH /albums/:id requests.
// Updates an album by its ID, allowing partial updates. Only provided fields are updated.
// Validates each provided field before updating. Returns the updated album as JSON with HTTP 200 status.
// Returns HTTP 400 if validation fails, or HTTP 404 if the album is not found.
func patchAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			var update Album
			if err := c.ShouldBindJSON(&update); err != nil {
				c.IndentedJSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid JSON",
					"details": err.Error(),
				})
				return
			}

			if update.Title != "" {
				if errMsg := validateTitle(update.Title, false); errMsg != "" {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
					return
				}
				albums[i].Title = update.Title
			}

			if update.Artist != "" {
				if errMsg := validateArtist(update.Artist, false); errMsg != "" {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
					return
				}
				albums[i].Artist = update.Artist
			}

			if update.Price > 0 {
				if errMsg := validatePrice(update.Price, false); errMsg != "" {
					c.IndentedJSON(http.StatusBadRequest, gin.H{"error": errMsg})
					return
				}
				albums[i].Price = update.Price
			}

			c.IndentedJSON(http.StatusOK, albums[i])
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
