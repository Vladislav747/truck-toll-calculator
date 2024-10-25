package grpc

import (
	"context"
	"github.com/Vladislav747/truck-toll-calculator/aggregator/service"
	"github.com/Vladislav747/truck-toll-calculator/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc service.Aggregator
}

func NewAggregatorGRPCServer(svc service.Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return nil, s.svc.AggregateDistance(distance)
}
