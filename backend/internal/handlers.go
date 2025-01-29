package internal

import (
	"encoding/json"
	"net/http"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	// Установка CORS-заголовков вручную (если не используешь middleware)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var request struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Здесь должна быть логика проверки кода
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Email verified successfully!"}`))
}
