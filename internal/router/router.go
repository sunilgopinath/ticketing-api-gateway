package router

import (
	"net/http"

	"github.com/sunilgopinath/ticketingapigateway/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/events", otelhttp.NewHandler(http.HandlerFunc(handlers.BrowseEventsHandler), "BrowseEvents"))
	mux.Handle("/purchase", otelhttp.NewHandler(http.HandlerFunc(handlers.PurchaseTicketHandler), "PurchaseTicket"))
	mux.Handle("/bookings", otelhttp.NewHandler(http.HandlerFunc(handlers.ViewBookingsHandler), "ViewBookings"))
	return mux
}
