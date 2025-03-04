package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func PurchaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("ticketing-api-gateway")
	ctx, span := tracer.Start(ctx, "PurchaseTicketHandler")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()

	// Generate cache key from request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Failed to read request body",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	cacheKey := generateCacheKeyFromBody(body, "purchase_")

	// Optional: Check cache (only if purchase is idempotent)
	cached, err := cache.GetCache(ctx, cacheKey)
	if err == nil && cached != "" {
		logger.Log.Info("Cache hit for /purchase",
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("client_ip", r.RemoteAddr),
			zap.String("cache_key", cacheKey),
		)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	logger.Log.Info("Cache miss for /purchase, processing request",
		zap.String("trace_id", traceID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", r.RemoteAddr),
		zap.String("cache_key", cacheKey),
	)

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

	// Stub response (replace with real purchase logic later)
	response := Response{Message: "Ticket purchase successful (stub)"}
	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.Log.Error("Failed to encode JSON response",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	// Optional: Cache response (if idempotent)
	if cacheErr := cache.SetCache(ctx, cacheKey, string(respBytes), 30*time.Second); cacheErr != nil {
		logger.Log.Warn("Failed to store response in cache",
			zap.String("trace_id", traceID),
			zap.String("cache_key", cacheKey),
			zap.Error(cacheErr),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// generateCacheKeyFromBody creates a cache key from POST body
func generateCacheKeyFromBody(body []byte, prefix string) string {
	hash := sha256.Sum256(body)
	return prefix + hex.EncodeToString(hash[:8])
}
