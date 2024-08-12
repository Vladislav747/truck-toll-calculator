package grpc

import (
	"github.com/Vladislav747/truck-toll-calculator/invoicer/service"
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

func (s *GRPCAggregatorServer) AggregateDistance(req types.AggregateRequest) error {
	distance := types.Distance{
		OBUID: int(req.ObuId),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return s.svc.AggregateDistance(distance)
}
