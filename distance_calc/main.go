package main

import (
	"fmt"
	consumer "github.com/Vladislav747/truck-toll-calculator/distance_calc/consumer"
	service "github.com/Vladislav747/truck-toll-calculator/distance_calc/service"
	"log"
)

var kafkaTopic string = "obuData"

//Transport (HTTP< GRPC) -> attcah business logic

func main() {
	var (
		err error
		svc service.CalculatorServicer
	)
	svc = service.NewCalculatorService()
	kafkaConsumer, err := consumer.NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("calc started")
}
