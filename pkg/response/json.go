package response

import (
	"encoding/json"
	"net/http"
)

// JSON writes a JSON response
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// Error writes a JSON error response
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, map[string]string{"error": message})
}

// Success writes a JSON success response
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, data)
}

// Created writes a JSON created response
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}


