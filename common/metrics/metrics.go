package metric

import "github.com/prometheus/client_golang/prometheus"

var Requests = prometheus.NewCounter(prometheus.CounterOpts{Name: "requests_total"})

func InitMetrics() { prometheus.MustRegister(Requests) }
