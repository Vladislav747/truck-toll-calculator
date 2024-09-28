package client

import (
	"context"
	"github.com/Vladislav747/truck-toll-calculator/types"
)

type Client interface {
	Aggregate(ctx context.Context, client *types.AggregateRequest) error
}
