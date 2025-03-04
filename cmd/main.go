package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"github.com/sunilgopinath/ticketingapigateway/internal/router"
	"github.com/sunilgopinath/ticketingapigateway/internal/tracing"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("API Gateway is starting...", zap.String("port", "8080"))

	// Initialize OpenTelemetry tracing
	shutdown := tracing.InitTracer()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

		if err := shutdown(ctx); err != nil {
			logger.Log.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}()

	routes := router.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", routes))
}
