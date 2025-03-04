package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Response struct {
	Message string `json:"message"`
}

func BrowseEventsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("ticketing-api-gateway") // Get OpenTelemetry tracer

	// Start a span for tracing
	ctx, span := tracer.Start(ctx, "BrowseEventsHandler")
	defer span.End()

	// Get real trace ID
	traceID := span.SpanContext().TraceID().String()

	cacheKey := generateCacheKey(r, "view_events_") // Reuse with unique prefix

	// Try retrieving from cache
	cached, err := cache.GetCache(ctx, cacheKey)
	if err == nil && cached != "" {
		logger.Log.Info("Cache hit for /events",
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

	// Cache miss - log and process request
	logger.Log.Info("Cache miss for /events, processing request",
		zap.String("trace_id", traceID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", r.RemoteAddr),
		zap.String("cache_key", cacheKey),
	)

	// Validate request method
	if r.Method != http.MethodGet {
		err := errors.New("invalid request method")
		logger.Log.Error("Invalid method for /events",
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("client_ip", r.RemoteAddr),
			zap.Error(err),
		)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Generate response
	response := Response{Message: fmt.Sprintf("List of events for query: %s", r.URL.RawQuery)}
	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.Log.Error("Failed to encode JSON response",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	// Store response in cache for 30 seconds
	if cacheErr := cache.SetCache(ctx, cacheKey, string(respBytes), 30*time.Second); cacheErr != nil {
		logger.Log.Warn("Failed to store response in cache",
			zap.String("trace_id", traceID),
			zap.String("cache_key", cacheKey),
			zap.Error(cacheErr),
		)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
