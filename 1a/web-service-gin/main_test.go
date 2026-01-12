package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupRouter creates a test router with all routes.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.PATCH("/albums/:id", patchAlbumByID)
	router.GET("/", healthCheck)
	return router
}

// resetAlbums resets albums to initial state for testing.
func resetAlbums() {
	albums = []Album{
		{ID: "550e8400-e29b-41d4-a716-446655440001", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "550e8400-e29b-41d4-a716-446655440002", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "550e8400-e29b-41d4-a716-446655440003", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
}

// TestHealthCheck tests the health check endpoint.
// Verifies that GET / returns HTTP 200 status.
func TestHealthCheck(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

// TestGetAlbums tests the GET /albums endpoint.
// Verifies that it returns HTTP 200 and at least 3 albums.
func TestGetAlbums(t *testing.T) {
	resetAlbums()
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/albums", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var albums []Album
	err := json.Unmarshal(w.Body.Bytes(), &albums)
	if err != nil {
		return
	}

	if len(albums) < 3 {
		t.Errorf("Expected at least 3 albums, got %d", len(albums))
	}
}

// TestGetAlbumByID tests the GET /albums/:id endpoint.
// Verifies successful retrieval returns HTTP 200, and non-existent ID returns HTTP 404.
func TestGetAlbumByID(t *testing.T) {
	resetAlbums()
	router := setupRouter()

	// Test existing album
	req, _ := http.NewRequest("GET", "/albums/550e8400-e29b-41d4-a716-446655440001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test non-existent album
	req, _ = http.NewRequest("GET", "/albums/not-found", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TestPostAlbums tests the POST /albums endpoint.
// Verifies that valid input creates an album with auto-generated ID (HTTP 201),
// and invalid input returns HTTP 400.
func TestPostAlbums(t *testing.T) {
	resetAlbums()
	router := setupRouter()

	// Test valid album creation
	body := `{"title": "Kind of Blue", "artist": "Miles Davis", "price": 49.99}`
	req, _ := http.NewRequest("POST", "/albums", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	var album Album
	err := json.Unmarshal(w.Body.Bytes(), &album)
	if err != nil {
		return
	}

	if album.ID == "" {
		t.Error("ID should be generated")
	}
	if album.Title != "Kind of Blue" {
		t.Errorf("Expected 'Kind of Blue', got '%s'", album.Title)
	}

	// Test invalid input
	invalidBody := `{"title": "A"}`
	req, _ = http.NewRequest("POST", "/albums", bytes.NewBufferString(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected 400 for invalid input, got %d", w.Code)
	}
}

// TestDeleteAlbumByID tests the DELETE /albums/:id endpoint.
// Verifies that deletion returns HTTP 200 and reduces album count,
// and non-existent ID returns HTTP 404.
func TestDeleteAlbumByID(t *testing.T) {
	resetAlbums()
	router := setupRouter()

	initialCount := len(albums)

	req, _ := http.NewRequest("DELETE", "/albums/550e8400-e29b-41d4-a716-446655440001", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	if len(albums) != initialCount-1 {
		t.Errorf("Expected %d albums, got %d", initialCount-1, len(albums))
	}

	// Test deleting non-existent album
	req, _ = http.NewRequest("DELETE", "/albums/not-found", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

// TestPatchAlbumByID tests the PATCH /albums/:id endpoint.
// Verifies that partial updates work correctly (HTTP 200),
// ID remains unchanged, and non-existent ID returns HTTP 404.
func TestPatchAlbumByID(t *testing.T) {
	resetAlbums()
	router := setupRouter()

	// Test updating title
	body := `{"title": "Updated Title"}`
	req, _ := http.NewRequest("PATCH", "/albums/550e8400-e29b-41d4-a716-446655440001", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var album Album
	err := json.Unmarshal(w.Body.Bytes(), &album)
	if err != nil {
		return
	}

	if album.Title != "Updated Title" {
		t.Errorf("Expected 'Updated Title', got '%s'", album.Title)
	}
	if album.ID != "550e8400-e29b-41d4-a716-446655440001" {
		t.Error("ID should not change")
	}

	// Test updating non-existent album
	req, _ = http.NewRequest("PATCH", "/albums/not-found", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}
