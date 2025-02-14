package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"
	"main/pkg/database"
)

// VerifyEmail verifies the email using a verification code
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decode request body
	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Here should be the logic to verify the code
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Email verified successfully!"}`))
}

// AddToCart adds an item to the cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	// Decode request body
	var req struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate input data
	if req.UserID == 0 || req.ProductID == 0 || req.Quantity <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user_id, product_id, or quantity"})
		return
	}

	// Check if the product exists in the database
	var product database.Product
	if err := database.DB.Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Product with ID %d not found", req.ProductID)})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch product details"})
		return
	}

	// Add the item to the cart
	cart := database.Cart{
		UserID:    uint(req.UserID),
		ProductID: uint(req.ProductID),
		Quantity:  req.Quantity,
	}
	if err := database.DB.Create(&cart).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to add item to cart"})
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart successfully"})
}

// InitiatePayment initiates the payment process
func InitiatePayment(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Decode request body
	var req struct {
		UserID int `json:"user_id"`
		CartID int `json:"cart_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validate input data
	if req.UserID == 0 || req.CartID == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid user_id or cart_id"})
		return
	}

	// Validate that the cart exists
	var cart database.Cart
	if err := database.DB.Where("cart_id = ?", req.CartID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Cart with ID %d not found", req.CartID)})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cart"})
		return
	}

	// Create a new transaction record
	transaction := database.Transaction{
		UserID: uint(req.UserID),
		CartID: uint(req.CartID),
		Status: "Pending Payment",
	}

	// Use raw SQL to insert and retrieve the generated transaction_id
	query := `
        INSERT INTO transactions (user_id, cart_id, status)
        VALUES ($1, $2, $3)
        RETURNING transaction_id
    `
	if err := database.DB.Raw(query, transaction.UserID, transaction.CartID, transaction.Status).Scan(&transaction).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create transaction"})
		return
	}

	// Validate TransactionID
	if transaction.TransactionID == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "TransactionID is empty after creation"})
		return
	}
	log.Printf("Transaction created with ID: %s", transaction.TransactionID)

	// Fetch cart items from the database
	var cartItems []database.Cart
	if err := database.DB.Where("cart_id = ?", req.CartID).Find(&cartItems).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch cart items"})
		return
	}

	// Prepare the payload dynamically
	var cartItemsPayload []map[string]interface{}
	for _, item := range cartItems {
		var product database.Product
		if err := database.DB.Where("product_id = ?", item.ProductID).First(&product).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Product with ID %d not found", item.ProductID)})
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch product details"})
			return
		}
		cartItemsPayload = append(cartItemsPayload, map[string]interface{}{
			"id":    item.ProductID,
			"name":  product.Name,
			"price": product.Price,
		})
	}

	// Fetch user details from the database
	var user database.User
	if err := database.DB.Where("user_id = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("User with ID %d not found", req.UserID)})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch user details"})
		return
	}

	payload := map[string]interface{}{
		"cartItems": cartItemsPayload,
		"customer": map[string]interface{}{
			"id":    req.UserID,
			"name":  user.UserName,
			"email": user.UserEmail,
		},
	}

	// Log the payload for debugging
	log.Printf("Payload sent to microservice: %+v", payload)

	// Send the POST request to the microservice
	microserviceURL := "http://localhost:8081/payment"
	resp, err := http.Post(microserviceURL, "application/json", encodeJSON(payload))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to process payment"})
		return
	}
	defer resp.Body.Close()

	// Parse the microservice response
	var microserviceResponse map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&microserviceResponse); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to parse microservice response"})
		return
	}

	log.Printf("Microservice response: %+v", microserviceResponse)

	// Validate microservice response
	success, exists := microserviceResponse["success"]
	if !exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "'success' key missing in microservice response"})
		return
	}

	// Update the transaction status based on the microservice response
	status := "Declined"
	if success {
		status = "Paid"
	}
	if err := database.DB.Model(&database.Transaction{}).
		Where("transaction_id = ?", transaction.TransactionID).
		Update("status", status).Error; err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update transaction status"})
		return
	}

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Payment " + status,
	})
}

// Helper function to encode JSON payload
func encodeJSON(data interface{}) *bytes.Buffer {
	payload, _ := json.Marshal(data)
	return bytes.NewBuffer(payload)
}
