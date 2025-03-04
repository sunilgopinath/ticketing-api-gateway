package handlers

import (
	"encoding/json"
	"net/http"
)

func PurchaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	response := Response{Message: "Ticket purchase successful (stub)"}
	json.NewEncoder(w).Encode(response)
}