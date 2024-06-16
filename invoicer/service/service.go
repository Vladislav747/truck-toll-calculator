package service

import (
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
	DistanceSum(int) (float64, error)
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
	Invoice(int) (*types.Invoice, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage:", distance)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) DistanceSum(obuId int) (float64, error) {
	fmt.Println("calculating sum distance in the storage - obuId: ", obuId)
	return i.store.Get(obuId)
}

func (i *InvoiceAggregator) CalculateInvoice(obuId int) (*types.Invoice, error) {
	fmt.Println("calculating invoice - obuId: ", obuId)
	dist, err := i.store.Get(obuId)
	if err != nil {
		logrus.Errorf("obu id (%d) not found", obuId)
		return nil, fmt.Errorf("obu id (%d) not found", obuId)
	}

	inv := &types.Invoice{
		OBUID:         obuId,
		TotalDistance: dist,
		TotalAmount:   0,
	}
	return inv, nil
}
