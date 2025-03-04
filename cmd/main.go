package main

import (
	"log"
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/logger"
	"github.com/sunilgopinath/ticketingapigateway/internal/router"
	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	logger.Log.Info("API Gateway is starting...", zap.String("port", "8080"))

	routes := router.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", routes))
}
