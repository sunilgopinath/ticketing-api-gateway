package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func ViewBookingsHandler(w http.ResponseWriter, r *http.Request, instance string) {

	ctx := r.Context()
	tracer := otel.Tracer("ticketing-api-gateway")
	ctx, span := tracer.Start(ctx, "ViewBookingsHandler")
	defer span.End()

	traceID := span.SpanContext().TraceID().String()
	cacheKey := generateCacheKey(r, "view_bookings_")

	// Pass endpoint to cache functions
	cached, err := cache.GetCache(ctx, cacheKey, "/bookings", instance)
	if err == nil && cached != "" {
		logger.Log.Info("Cache hit for /bookings",
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

	logger.Log.Info("Cache miss for /bookings, processing request",
		zap.String("trace_id", traceID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("client_ip", r.RemoteAddr),
		zap.String("cache_key", cacheKey),
	)

	if r.Method != http.MethodGet {
		err := errors.New("invalid request method")
		logger.Log.Error("Invalid method for /bookings",
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("client_ip", r.RemoteAddr),
			zap.Error(err),
		)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	response := Response{Message: fmt.Sprintf("User bookings for query: %s", r.URL.RawQuery)}
	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.Log.Error("Failed to encode JSON response",
			zap.String("trace_id", traceID),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}

	if cacheErr := cache.SetCache(ctx, cacheKey, string(respBytes), 30*time.Second, "/bookings", instance); cacheErr != nil {
		logger.Log.Warn("Failed to store response in cache",
			zap.String("trace_id", traceID),
			zap.String("cache_key", cacheKey),
			zap.Error(cacheErr),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

// Modified generateCacheKey to accept a prefix
func generateCacheKey(r *http.Request, prefix string) string {
	queryParams := r.URL.Query()
	var keys []string
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryString []string
	for _, k := range keys {
		queryString = append(queryString, fmt.Sprintf("%s=%s", k, strings.Join(queryParams[k], ",")))
	}
	joinedParams := strings.Join(queryString, "&")

	hash := sha256.Sum256([]byte(joinedParams))
	return prefix + hex.EncodeToString(hash[:8])
}
