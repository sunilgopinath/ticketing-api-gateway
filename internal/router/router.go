package router

import (
	"net"
	"net/http"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/handlers"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

var limiter *redis_rate.Limiter

func rateLimit(next http.Handler, instance string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Log.Error("Failed to parse client IP", zap.Error(err))
			clientIP = r.RemoteAddr // Fallback, though unlikely needed
		}
		key := "ratelimit:" + clientIP + ":" + instance
		limit := redis_rate.Limit{
			Rate:   10,          // 10 requests
			Burst:  10,          // Allow burst up to 10
			Period: time.Minute, // Per minute
		}
		res, err := limiter.Allow(r.Context(), key, limit)
		if err != nil {
			logger.Log.Error("Failed to check rate limit", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if res.Allowed == 0 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func SetupRoutes(instance string) *http.ServeMux {
	if limiter == nil {
		limiter = redis_rate.NewLimiter(cache.RedisClient)
	}

	mux := http.NewServeMux()
	mux.Handle("/events", rateLimit(otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.BrowseEventsHandler(w, r, instance)
	}), "BrowseEvents"), instance))
	mux.Handle("/purchase", rateLimit(otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.PurchaseTicketHandler(w, r, instance)
	}), "PurchaseTicket"), instance))
	mux.Handle("/bookings", rateLimit(otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ViewBookingsHandler(w, r, instance)
	}), "ViewBookings"), instance))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
