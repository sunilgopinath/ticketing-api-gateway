package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"go.uber.org/zap"
)

func PurchaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	traceID := "random-trace-id-purchase-1234"
	if r.Method != http.MethodPost {
		err := errors.New("invalid request method")
		logger.Log.Error("Invalid method for /purchase",
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("client_ip", r.RemoteAddr),
			zap.Error(err),
		)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	logger.Log.Info("Processing /purchase request",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", r.RemoteAddr),
	)
	response := Response{Message: "Ticket purchase successful (stub)"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.Error("Failed to encode JSON response",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
