package main

import (
	"encoding/json"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/invoicer/middleware"
	"github.com/Vladislav747/truck-toll-calculator/invoicer/service"
	store2 "github.com/Vladislav747/truck-toll-calculator/invoicer/store"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"net/http"
)

func main() {

	var (
		store = store2.NewMemoryStore()
		svc   = service.NewInvoiceAggregator(store)
		mid   = middleware.NewLogMiddleware(svc)
	)
	makeHTTPTransport(":3002", mid)
}

func makeHTTPTransport(listenAddr string, svc service.Aggregator) {
	fmt.Println("HTTP Transport Listening on " + listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc service.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
