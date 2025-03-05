package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/cache"
	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"github.com/sunilgopinath/ticketingapigateway/internal/router"
	"github.com/sunilgopinath/ticketingapigateway/internal/tracing"
	"go.uber.org/zap"
)

func main() {
	port := flag.String("port", "8080", "Port to run the server on")
	instance := flag.String("instance", "gateway-1", "Instance ID for this server")
	flag.Parse()

	logger.InitLogger()
	logger.Log.Info("API Gateway is starting...", zap.String("port", *port), zap.String("instance", *instance))

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
	routes := router.SetupRoutes(*instance)
	log.Fatal(http.ListenAndServe(":"+*port, routes))
}
