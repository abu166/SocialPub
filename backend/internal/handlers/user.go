package handlers

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"main/internal/models"
	"main/internal/utils"
	"main/pkg/database"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	maxRequestBodySize = 1024 * 1024 // 1MB
	maxQueryParamSize  = 1024        // 1KB
)

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Allow only POST method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := ResponseData{Status: "fail", Message: "Only POST method is allowed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Limit request body size
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	// Check if request body is empty
	if r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "Request body cannot be empty"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse JSON body and detect unexpected fields
	var requestData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Disallow extra/unknown fields
	err := decoder.Decode(&requestData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "Invalid JSON or unexpected fields"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check "message" key
	message, ok := requestData["message"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "Message field is required"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate "message" value
	messageStr, isString := message.(string)
	if !isString {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "Message field must be a string"}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Printf("Received message: %s\n", messageStr)

	// Send success response
	response := ResponseData{Status: "success", Message: "Data successfully received"}
	json.NewEncoder(w).Encode(response)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Allow only GET method
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Check query string size
	queryString := r.URL.RawQuery
	if len(queryString) > maxQueryParamSize {
		http.Error(w, "Query string size exceeds limit", http.StatusRequestEntityTooLarge)
		return
	}

	// Parse and log query parameters
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		fmt.Println("No query parameters provided")
	} else {
		fmt.Println("Query Parameters:", queryParams)
	}

	// Validate "message" query parameter
	msg := queryParams.Get("message")
	if strings.TrimSpace(msg) == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "Missing 'message' query parameter"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check message length
	if len(msg) > 256 {
		w.WriteHeader(http.StatusBadRequest)
		response := ResponseData{Status: "fail", Message: "'message' parameter is too long"}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println("Query parameter 'message':", msg)

	// Send success response
	response := ResponseData{Status: "success", Message: "GET request received"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		CreateUser(w, r)
	case http.MethodGet:
		GetAllUsers(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	logger := utils.Log.WithField("handler", "createUser")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.HandleError(w, utils.ErrInvalidInput, http.StatusBadRequest, logger)
		return
	}

	// Set created_at and updated_at timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Validate user fields
	if err := utils.ValidateUser(user); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, logger)
		return
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		if utils.IsDuplicateEmailError(err) {
			utils.HandleError(w, utils.ErrDuplicateEmail, http.StatusConflict, logger)
			return
		}
		utils.HandleError(w, utils.ErrDatabaseOperation, http.StatusInternalServerError, logger)
		return
	}

	// Send response
	utils.SendJSONResponse(w, http.StatusCreated, models.ResponseData{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	})
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filterField := r.URL.Query().Get("filter_field")
	filterValue := r.URL.Query().Get("filter_value")
	filterOperator := r.URL.Query().Get("filter_operator")
	sortField := r.URL.Query().Get("sort_field")
	sortDir := r.URL.Query().Get("sort_dir")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	itemsPerPage := 5

	query := database.DB.Model(&models.User{})

	if filterValue != "" {
		query = applyFilter(query, filterField, filterValue, filterOperator)
	}

	var totalItems int64
	query.Count(&totalItems)

	query = applySorting(query, sortField, sortDir)

	offset := (page - 1) * itemsPerPage
	query = query.Offset(offset).Limit(itemsPerPage)

	var users []models.User
	if err := query.Find(&users).Error; err != nil {
		utils.SendErrorResponse(w, "Could not retrieve users", http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(itemsPerPage)))

	utils.SendJSONResponse(w, http.StatusOK, models.PaginatedResponse{
		Status:      "success",
		Message:     "Users retrieved successfully",
		Data:        users,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
	})
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var updateData struct {
		UserID    uint   `json:"user_id"`
		UserName  string `json:"user_name,omitempty"`
		UserEmail string `json:"user_email,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.First(&user, updateData.UserID).Error; err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	if updateData.UserName != "" {
		user.UserName = updateData.UserName
	}
	if updateData.UserEmail != "" {
		if !utils.IsValidEmail(updateData.UserEmail) {
			utils.SendErrorResponse(w, "Invalid email format", http.StatusBadRequest)
			return
		}
		user.UserEmail = updateData.UserEmail
	}

	user.UpdatedAt = time.Now()

	if err := database.DB.Save(&user).Error; err != nil {
		utils.SendErrorResponse(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, models.ResponseData{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var deleteID uint
	var err error

	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		var id uint64
		id, err = strconv.ParseUint(idStr, 10, 32)
		deleteID = uint(id)
	} else {
		var requestBody map[string]uint
		err = json.NewDecoder(r.Body).Decode(&requestBody)
		if err == nil {
			deleteID = requestBody["user_id"]
		}
	}

	if err != nil || deleteID == 0 {
		utils.SendErrorResponse(w, "Invalid or missing user ID", http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&models.User{}, deleteID)
	if result.Error != nil {
		utils.SendErrorResponse(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		utils.SendErrorResponse(w, "No user found with the given ID", http.StatusNotFound)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, models.ResponseData{
		Status:  "success",
		Message: "User deleted successfully",
	})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.SendErrorResponse(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.SendErrorResponse(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.First(&user, uint(id)).Error; err != nil {
		utils.SendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, models.ResponseData{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// Helper functions for filtering and sorting
func applyFilter(query *gorm.DB, field, value, operator string) *gorm.DB {
	switch field {
	case "name":
		return applyStringFilter(query, "user_name", value, operator)
	case "email":
		return applyStringFilter(query, "user_email", value, operator)
	case "date":
		return applyDateFilter(query, "created_at", value, operator)
	default:
		return query
	}
}

func applyStringFilter(query *gorm.DB, field, value, operator string) *gorm.DB {
	switch operator {
	case "contains":
		return query.Where(field+" ILIKE ?", "%"+value+"%")
	case "equals":
		return query.Where(field+" = ?", value)
	case "startsWith":
		return query.Where(field+" ILIKE ?", value+"%")
	case "endsWith":
		return query.Where(field+" ILIKE ?", "%"+value)
	default:
		return query
	}
}

func applyDateFilter(query *gorm.DB, field, value, operator string) *gorm.DB {
	switch operator {
	case "equals":
		return query.Where(field+" BETWEEN ? AND ?", value, value+"T23:59:59Z")
	case "before":
		return query.Where(field+" < ?", value)
	case "after":
		return query.Where(field+" > ?", value)
	default:
		return query
	}
}

func applySorting(query *gorm.DB, field, direction string) *gorm.DB {
	if field == "" {
		return query.Order("user_id asc")
	}

	if direction != "desc" {
		direction = "asc"
	}

	return query.Order(fmt.Sprintf("%s %s", field, direction))
}
