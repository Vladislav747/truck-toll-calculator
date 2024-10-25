package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/client"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address of the HTTP server")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:3000", "the listen address of the aggregator service")

	flag.Parse()
	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr) //endpoint of the aggregator service
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port: %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{client: c}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("hitting the get invoice inside the gateway ")
	inv, err := h.client.GetInvoice(context.Background(), 123)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)

}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
			}).Info()
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
