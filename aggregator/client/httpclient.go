package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	invReq := &types.GetInvoiceRequest{
		ObuID: int32(id),
	}

	b, err := json.Marshal(invReq)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s?obu=%d", c.Endpoint, "invoice", id)
	logrus.Debugln("requesting get invoice ->", endpoint)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		logrus.Error(err, "Send query to /invoice/aggregate")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("The service responded with non 200 status: %d", resp.StatusCode)
		return nil, fmt.Errorf("the service responded with non 200 status: %d", resp.StatusCode)
	}

	var inv types.Invoice
	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return &inv, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	b, err := json.Marshal(aggReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint+"/invoice", bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		logrus.Error(err, "Send query to /invoice/aggregate")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("The service responded with non 200 status: %d", resp.StatusCode)
		return fmt.Errorf("the service responded with non 200 status: %d", resp.StatusCode)
	}
	resp.Body.Close()
	return nil
}
