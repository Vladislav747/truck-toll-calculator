package middleware

import (
	"github.com/Vladislav747/truck-toll-calculator/data_receiver/producer"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/sirupsen/logrus"
	"time"
)

type LoggingMiddleware struct {
	next producer.DataProducer
}

func NewLoggingMiddleware(next producer.DataProducer) *LoggingMiddleware {
	return &LoggingMiddleware{next: next}
}

func (l *LoggingMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"long":  data.Long,
			"lat":   data.Lat,
			"took":  time.Since(start),
		}).Info("producting to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
