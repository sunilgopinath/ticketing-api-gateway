package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/router"
)

func main() {
	fmt.Println("API Gateway is running on port 8080...")
	routes := router.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", routes))
}
