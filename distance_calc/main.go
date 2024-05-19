package main

import (
	"fmt"
	consumer2 "github.com/Vladislav747/truck-toll-calculator/distance_calc/consumer"
	"log"
)

type DistanceCalculator struct {
	consumer consumer2.DataConsumer
}

var kafkaTopic string = "obudata"

func main() {
	kafkaConsumer, err := consumer2.NewKafkaConsumer(kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	fmt.Println("calc started")
}
