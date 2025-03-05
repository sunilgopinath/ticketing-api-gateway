package router

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sunilgopinath/ticketingapigateway/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func SetupRoutes(instance string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/events", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.BrowseEventsHandler(w, r, instance)
	}), "BrowseEvents"))
	mux.Handle("/purchase", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.PurchaseTicketHandler(w, r, instance)
	}), "PurchaseTicket"))
	mux.Handle("/bookings", otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ViewBookingsHandler(w, r, instance)
	}), "ViewBookings"))
	mux.Handle("/metrics", promhttp.Handler())
	return mux
}
