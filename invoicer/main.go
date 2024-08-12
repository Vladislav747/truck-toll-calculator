package main

import (
	"encoding/json"
	"flag"
	"fmt"
	grpc2 "github.com/Vladislav747/truck-toll-calculator/invoicer/grpc"
	"github.com/Vladislav747/truck-toll-calculator/invoicer/middleware"
	"github.com/Vladislav747/truck-toll-calculator/invoicer/service"
	store2 "github.com/Vladislav747/truck-toll-calculator/invoicer/store"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"strconv"
)

func main() {

	httpListenAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen address of the GRPC server")

	flag.Parse()
	var (
		store = store2.NewMemoryStore()
		svc   = service.NewInvoiceAggregator(store)
		mid   = middleware.NewLogMiddleware(svc)
	)
	go makeGRPCTransport(*grpcListenAddr, mid)
	makeHTTPTransport(*httpListenAddr, mid)
}

func makeHTTPTransport(listenAddr string, svc service.Aggregator) {
	fmt.Println("HTTP Transport Listening on " + listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc service.Aggregator) error {
	fmt.Println("gRPC Transport Listening on " + listenAddr)
	//Make a TCP Listeners
	ln, err := net.Listen("TCP", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	//Make a new GRPC native Server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register (OUR) GRPC server implementation to the GRPC Implemetation
	types.RegisterAggregatorServer(server, grpc2.NewAggregatorGRPCServer(svc))

	return server.Serve(ln)
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

func handleInvoice(svc service.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
