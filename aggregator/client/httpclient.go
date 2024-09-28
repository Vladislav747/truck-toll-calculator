package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err, "Send query to /invoice/aggregate")
	}
	return nil
}
