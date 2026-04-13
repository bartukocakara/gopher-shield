package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	CBTransitions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gophershield_cb_transitions_total",
			Help: "Circuit breaker state changes",
		},
		[]string{"state"},
	)
)

func init() {
	prometheus.MustRegister(CBTransitions)
}