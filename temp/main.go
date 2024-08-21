package main

import (
	"context"
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":3001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	aggrConn := types.NewAggregatorClient(conn)
	fmt.Println(aggrConn)

	if _, err := aggrConn.Aggregate(context.Background(), &types.AggregateRequest{
		ObuId: 1,
		Value: 56.60,
		Unix:  time.Now().UnixNano(),
	}); err != nil {

	}
}
