package client

import (
	"bytes"
	"encoding/json"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"net/http"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", c.Endpoint+"/invoice/aggregate", bytes.NewReader(b))
	if err != nil {
		return err
	}
	http.DefaultClient.Do(req)
	return nil
}
