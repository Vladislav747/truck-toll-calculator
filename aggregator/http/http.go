package http

import (
	"encoding/json"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/service"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type HTTPMetricsHandler struct {
	reqCounter prometheus.Counter
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricsHandler {
	reqCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "http_request_total",
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "req_counter"),
		Help:      "Total number of HTTP requests",
	})
	return &HTTPMetricsHandler{
		reqCounter: reqCounter,
	}
}

func (h *HTTPMetricsHandler) instrument(next http.Handler) http.HandlerFunc {
	h.reqCounter.Inc()
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}

func HandleAggregate(svc service.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "method POST is not allowed"})
			return
		}
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

func HandleInvoice(svc service.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "method GET is not allowed"})
			return
		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}
		fmt.Println(values, "values")
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			logrus.Println("error converting to int", err)
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}
		distance, err := svc.CalculateInvoice(obuId)
		if err != nil {
			logrus.Println("error calculating", err)
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error calculating"})
		}
		writeJSON(w, http.StatusOK, map[string]any{"distance": distance})
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
