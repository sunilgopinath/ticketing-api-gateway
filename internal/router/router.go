package router

import (
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", handlers.BrowseEventsHandler)
	mux.HandleFunc("/purchase", handlers.PurchaseTicketHandler)
	mux.HandleFunc("/bookings", handlers.ViewBookingsHandler)
	return mux
}
