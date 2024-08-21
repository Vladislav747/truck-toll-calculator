package client

import (
	"github.com/Vladislav747/truck-toll-calculator/types"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)

	return &GRPCClient{Endpoint: endpoint, AggregatorClient: c}, nil
}
