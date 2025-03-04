package handlers

import (
	"encoding/json"
	"net/http"
)

func ViewBookingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	response := Response{Message: "User booksings (stub)"}
	json.NewEncoder(w).Encode(response)
}
