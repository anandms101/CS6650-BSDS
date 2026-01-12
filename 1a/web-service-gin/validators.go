package main

// validateTitle validates the title field and returns an error message if validation fails.
// If required is true, the title must be non-empty. The title must be between 2 and 100 characters.
// Returns an empty string if validation passes, otherwise returns an error message.
func validateTitle(title string, required bool) string {
	if required && title == "" {
		return "Title is required"
	}
	if title != "" && (len(title) < 2 || len(title) > 100) {
		return "Title must be between 2 and 100 characters"
	}
	return ""
}

// validateArtist validates the artist field and returns an error message if validation fails.
// If required is true, the artist must be non-empty. The artist must be between 2 and 100 characters.
// Returns an empty string if validation passes, otherwise returns an error message.
func validateArtist(artist string, required bool) string {
	if required && artist == "" {
		return "Artist is required"
	}
	if artist != "" && (len(artist) < 2 || len(artist) > 100) {
		return "Artist must be between 2 and 100 characters"
	}
	return ""
}

// validatePrice validates the price field and returns an error message if validation fails.
// If required is true, the price must be greater than 0. Price cannot be negative.
// Returns an empty string if validation passes, otherwise returns an error message.
func validatePrice(price float64, required bool) string {
	if required && price <= 0 {
		return "Price is required and must be greater than 0"
	}
	if price < 0 {
		return "Price must be greater than or equal to 0"
	}
	return ""
}
