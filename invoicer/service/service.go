package service

import (
	"fmt"
	"github.com/Vladislav747/truck-toll-calculator/types"
)

type Aggregator interface {
	AggregateDistance(distance types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
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
