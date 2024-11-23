package metrics

import (
	"github.com/jatin297/retoenfa/dto"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

// metrics collectors
var (
	httpRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: string(dto.TotalRequests),
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    string(dto.RequestDuration),
			Help:    "Histogram of response time for handler in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	regexSizeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: string(dto.RegexSize),
			Help: "Size of the regular expression being processed",
		},
		[]string{"path"},
	)

	transitionTableSize = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: string(dto.EnfaTransitionTableSize),
			Help: "Number of transitions in the ENFA",
		},
		[]string{"path"},
	)

	regexProcessedCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: string(dto.RegexProcessedTotal),
			Help: "Total number of successfully processed regular expressions",
		},
		[]string{"path"},
	)

	processingTimeBySize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    string(dto.ProcessingTimeByRegexSize),
			Help:    "Histogram of processing time bucketed by regex size",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"path"},
	)
)

func RecordMetrics(method, path string, statusCode int, start time.Time, regularExpression dto.RegularExpression, eNFA dto.ENFAResponse) {
	status := http.StatusText(statusCode)
	// Increment request count
	httpRequestCount.WithLabelValues(method, path, status).Inc()

	// Observe request duration
	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)

	recordRegexSize(path, len(regularExpression.RE))
	recordTransitionTableSize(path, eNFA.TransitionTableSize)
	recordRegexProcessed(path)
	recordProcessingTimeBySize(path, len(regularExpression.RE), duration)

}

func RecordMetricForHttp(method, path string, statusCode int, start time.Time) {
	status := http.StatusText(statusCode)
	// Increment request count
	httpRequestCount.WithLabelValues(method, path, status).Inc()

	// Observe request duration
	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues(method, path, status).Observe(duration)
}

func recordRegexSize(path string, size int) {
	regexSizeGauge.WithLabelValues(path).Set(float64(size))
}

func recordTransitionTableSize(path string, size int) {
	transitionTableSize.WithLabelValues(path).Set(float64(size))
}

func recordRegexProcessed(path string) {
	regexProcessedCount.WithLabelValues(path).Inc()
}

func recordProcessingTimeBySize(path string, size int, duration float64) {
	processingTimeBySize.WithLabelValues(path).Observe(duration)
}

func init() {
	// Register metrics
	prometheus.MustRegister(httpRequestCount)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(regexSizeGauge)
	prometheus.MustRegister(transitionTableSize)
	prometheus.MustRegister(regexProcessedCount)
	prometheus.MustRegister(processingTimeBySize)
}
