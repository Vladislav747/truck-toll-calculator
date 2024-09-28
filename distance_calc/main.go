package main

import (
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/client"
	consumer "github.com/Vladislav747/truck-toll-calculator/distance_calc/consumer"
	"github.com/Vladislav747/truck-toll-calculator/distance_calc/middleware"
	service "github.com/Vladislav747/truck-toll-calculator/distance_calc/service"
	"log"
)

const (
	kafkaTopic         = "obuData"
	aggregatorEndpoint = "http://127.0.0.1:3002/aggregate"
)

//Transport (HTTP< GRPC) -> attcah business logic

func main() {
	var (
		err error
		svc service.CalculatorServicer
	)
	svc = service.NewCalculatorService()
	svc = middleware.NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	//grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	//if err != nil {
	//	log.Fatal(err)
	//}

	kafkaConsumer, err := consumer.NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("calc started")
}
