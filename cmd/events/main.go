package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/handlers"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"github.com/sunilgopinath/ticketingapigateway/internal/tracing"
	"go.uber.org/zap"
)

func main() {

	port := flag.String("port", "8081", "Port for events service")
	instance := flag.String("instance", "gateway-1", "Instance ID for this server")
	flag.Parse()

	logger.InitLogger()

	cache.InitRedis()

	shutdown, err := tracing.InitTracer()
	if err != nil {
		logger.Log.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		ctx := context.Background()
		if err := shutdown(ctx); err != nil {
			logger.Log.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	// Pass instance to router
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		handlers.BrowseEventsHandler(w, r, *instance)
	})
	logger.Log.Info("Events service starting on", zap.String("port", *port))
	log.Fatal(http.ListenAndServe(":"+*port, mux))
}
