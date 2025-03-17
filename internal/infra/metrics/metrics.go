package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	URLShortened = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "urls_shortened_total",
		Help: "Total number of shortened URLs",
	})

	URLRedirects = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "url_redirects_total",
		Help: "Total number of URL redirects",
	}, []string{"url_id"})

	ResponseTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "HTTP response time in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "method"})
)

func init() {
	prometheus.MustRegister(URLShortened)
	prometheus.MustRegister(URLRedirects)
	prometheus.MustRegister(ResponseTime)
}
