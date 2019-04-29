package state

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

const MetricsSubsystem = "state"

type Metrics struct {
	// Time between BeginBlock and EndBlock.
	BlockProcessingTime metrics.Histogram

	// Time on recheck
	RecheckTime metrics.Histogram

	// App hash conflict error
	AppHashConflict metrics.Counter
}

func PrometheusMetrics(namespace string) *Metrics {
	return &Metrics{
		BlockProcessingTime: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "block_processing_time",
			Help:      "Time between BeginBlock and EndBlock in ms.",
			Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
		}, []string{}),
		RecheckTime: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "recheck_time",
			Help:      "Time cost on recheck in ms.",
			Buckets:   stdprometheus.LinearBuckets(1, 10, 10),
		}, []string{}),
		AppHashConflict: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: MetricsSubsystem,
			Name:      "app_hash_conflict",
			Help:      "App hash conflict error",
		}, []string{"proposer", "height"}),
	}
}

func NopMetrics() *Metrics {
	return &Metrics{
		BlockProcessingTime: discard.NewHistogram(),
		RecheckTime: 	     discard.NewHistogram(),
		AppHashConflict:	 discard.NewCounter(),
	}
}
