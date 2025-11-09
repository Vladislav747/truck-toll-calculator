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
	"time"
)

type HTTPMetricsHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHTTPMetricsHandler(reqName string) *HTTPMetricsHandler {

	reqCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      fmt.Sprintf("http_%s", reqName),
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Help:      "Total number of HTTP requests",
	})

	reqLatency := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      fmt.Sprintf("http_%s", reqName),
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Help:      "Latency of HTTP requests",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	prometheus.MustRegister(reqCounter)
	prometheus.MustRegister(reqLatency)
	return &HTTPMetricsHandler{
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricsHandler) Instrument(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
			}).Info()
			h.reqLatency.Observe(time.Since(start).Seconds())
		}(time.Now())
		h.reqCounter.Inc()
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
