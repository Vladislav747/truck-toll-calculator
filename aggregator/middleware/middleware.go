package middleware

import (
	"github.com/Vladislav747/truck-toll-calculator/aggregator/service"
	"github.com/Vladislav747/truck-toll-calculator/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"time"
)

type MetricsMiddleware struct {
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           service.Aggregator
}

func NewMetricsMiddleware(next service.Aggregator) *MetricsMiddleware {
	errCounterAgg := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "calc_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "calc_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	prometheus.MustRegister(errCounterAgg)
	prometheus.MustRegister(errCounterCalc)
	prometheus.MustRegister(reqCounterAgg)
	prometheus.MustRegister(reqCounterCalc)
	prometheus.MustRegister(reqLatencyAgg)
	prometheus.MustRegister(reqLatencyCalc)
	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
	}
}

func (m *MetricsMiddleware) pushMetrics(start time.Time, err error) {

}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return err
}

func (m *MetricsMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	return inv, err
}

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
