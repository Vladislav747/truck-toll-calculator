package main

import (
	"flag"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/data_receiver/middleware"
	"github.com/Vladislav747/truck-toll-calculator/data_receiver/producer"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var kafkaTopic string = "obuData"
var defaultPort string = "30000"

func main() {
	// Загрузить файл env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading.env file")
	}
	flag.Parse()

	port := os.Getenv("DATA_RECEIVER_PORT")
	if port == "" {
		port = defaultPort
	}

	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWS)
	fmt.Println("data_receiver started on", port, "port")
	addr := ":" + port
	http.ListenAndServe(addr, nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  producer.DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p   producer.DataProducer
		err error
	)
	p, err = producer.NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	p = middleware.NewLoggingMiddleware(p)
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()

}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("client connected to ws")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		fmt.Printf("received OBU data from [%d] :: <lat %.2f, long %.2f> \n", data.OBUID, data.Lat, data.Long)
		dr.msgch <- data
		if err := dr.produceData(data); err != nil {
			log.Println("kafka produce error: ", err)
		}
	}
}
