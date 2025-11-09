package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/client"
	grpc2 "github.com/Vladislav747/truck-toll-calculator/aggregator/grpc"
	http2 "github.com/Vladislav747/truck-toll-calculator/aggregator/http"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/middleware"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/service"
	store2 "github.com/Vladislav747/truck-toll-calculator/aggregator/store"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading.env file")
	}
	flag.Parse()
	var (
		store          = makeStore()
		svc            = service.NewInvoiceAggregator(store)
		mid            = middleware.NewLogMiddleware(svc)
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
	)
	mid = middleware.NewMetricsMiddleware(svc)
	mid = middleware.NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, mid))
	}()
	time.Sleep(time.Second * 5)
	c, err := client.NewGRPCClient(grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to gRPC server", c)
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuId: 1,
		Value: 56.60,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	log.Fatal(makeHTTPTransport(httpListenAddr, mid))
}

func makeHTTPTransport(listenAddr string, svc service.Aggregator) error {
	aggMetricsHandler := http2.NewHTTPMetricsHandler("aggregator")
	invoiceMetricsHandler := http2.NewHTTPMetricsHandler("invoice")
	http.HandleFunc("/aggregate", aggMetricsHandler.Instrument(http2.HandleAggregate(svc)))
	http.HandleFunc("/invoice", invoiceMetricsHandler.Instrument(http2.HandleInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP Transport Listening on " + listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc service.Aggregator) error {
	fmt.Println("gRPC Transport Listening on " + listenAddr)
	//Make a TCP Listeners
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer func() {
		fmt.Println("stopping GRPC transport")
		ln.Close()
	}()
	//Make a new GRPC native Server
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register (OUR) GRPC server implementation to the GRPC Implemetation
	types.RegisterAggregatorServer(server, grpc2.NewAggregatorGRPCServer(svc))

	return server.Serve(ln)
}

func makeStore() service.Storer {
	storerType := os.Getenv("AGG_STORAGE_TYPE")
	switch storerType {
	case "memory":
		return store2.NewMemoryStore()
	default:
		log.Fatal("Unsupported storage type: %s", storerType)
		return nil
	}
}
