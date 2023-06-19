package unknownip

import (
	"sync"

	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// requestCount exports a prometheus metric that is incremented every time a query is seen by the unknownip plugin.
var requestCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "unknownip",
	Name:      "request_count_total",
	Help:      "Counter of requests made.",
}, []string{"server"})

// unknownIpCount exports a prometheus metric that is incremented every time a query return an IP not contained in the given prefixes.
var unknownIpCount = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: "unknownip",
	Name:      "unknown_count_total",
	Help:      "Counter of requests that returned an unknown IP.",
}, []string{"server"})

var once sync.Once
