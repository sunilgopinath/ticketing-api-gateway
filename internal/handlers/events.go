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

func BrowseEventsHandler(w http.ResponseWriter, r *http.Request, instance string) {
	ctx := r.Context()
	tracer := otel.Tracer("ticketing-api-gateway")
	ctx, span := tracer.Start(ctx, "BrowseEventsHandler")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	cacheKey := generateCacheKey(r, "browse_events_")

	cached, err := cache.GetCache(ctx, cacheKey, "/events", instance)
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

	logger.Log.Info("Cache miss for /events, processing request",
		zap.String("trace_id", traceID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", r.RemoteAddr),
		zap.String("cache_key", cacheKey),
	)

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

	if cacheErr := cache.SetCache(ctx, cacheKey, string(respBytes), 30*time.Second, "/events", instance); cacheErr != nil {
		logger.Log.Warn("Failed to store response in cache",
			zap.String("trace_id", traceID),
			zap.String("cache_key", cacheKey),
			zap.Error(cacheErr),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}
