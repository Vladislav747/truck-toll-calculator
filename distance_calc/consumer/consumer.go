package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/distance_calc/service"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type DataConsumer interface {
	ConsumeData()
}

type kafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService service.CalculatorServicer
}

func NewKafkaConsumer(topic string, svc service.CalculatorServicer) (*kafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &kafkaConsumer{
		consumer:    c,
		calcService: svc,
	}, nil
}

func (c *kafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *kafkaConsumer) Stop() {
	logrus.Info("kafka transport stopped")
	c.isRunning = false
}

func (c *kafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Calculation error %s", err)
			continue
		}
		fmt.Println(distance, "distance")
	}
}
