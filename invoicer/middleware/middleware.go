package middleware

import (
	"github.com/Vladislav747/truck-toll-calculator/invoicer/service"
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
			"func": "AggregateDistance",
		}).Info()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}
