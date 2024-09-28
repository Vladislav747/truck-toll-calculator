package middleware

import (
	"github.com/Vladislav747/truck-toll-calculator/aggregator/service"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LogMiddleware struct {
	next service.Aggregator
}

func NewLogMiddleware(next service.Aggregator) service.Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {

	defer func(start time.Time) {

		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"err":         err,
			"obuID":       obuId,
			"totalDist":   distance,
			"totalAmount": amount,
		}).Info("CalculateInvoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	return
}
