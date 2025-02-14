package internal

import (
	"encoding/json"
	"log"
	"main/pkg/database"
	"net/http"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract user ID from session or token (mocked here for simplicity)
	userID := 4 // Replace with dynamic user ID

	// Fetch cart items from the database
	var cartItems []database.Cart
	if err := database.DB.Where("user_id = ?", userID).Find(&cartItems).Error; err != nil {
		log.Printf("Error fetching cart items for user_id=%d: %v", userID, err)
		http.Error(w, "Failed to fetch cart items", http.StatusInternalServerError)
		return
	}

	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cartItems": cartItems,
	})
}
