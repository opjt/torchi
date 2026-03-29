package middle

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// duration 측정에서 제외할 경로 suffix
// - /ask: 사용자 반응 대기(최대 300s)로 레이턴시 왜곡
// - /sse: 장기 연결로 레이턴시 의미 없음
var skipDurationSuffixes = []string{"/ask", "/sse/notifications"}

type Metrics struct {
	requestsTotal    *prometheus.CounterVec
	requestDuration  *prometheus.HistogramVec
	requestsInFlight prometheus.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		requestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		}, []string{"method", "path", "status"}),

		requestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		}, []string{"method", "path"}),

		requestsInFlight: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Number of HTTP requests currently being processed.",
		}),
	}
}

func (m *Metrics) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.requestsInFlight.Inc()
			defer m.requestsInFlight.Dec()

			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(ww, r)

			// chi route pattern 사용으로 고카디널리티 방지
			// e.g. /api/v1/push/{token} (실제 토큰 값 대신)
			path := chi.RouteContext(r.Context()).RoutePattern()
			if path == "" {
				path = r.URL.Path
			}

			if path == "/metrics" {
				return
			}

			status := strconv.Itoa(ww.status)
			duration := time.Since(start).Seconds()

			m.requestsTotal.WithLabelValues(r.Method, path, status).Inc()

			skipDuration := false
			for _, suffix := range skipDurationSuffixes {
				if strings.HasSuffix(path, suffix) {
					skipDuration = true
					break
				}
			}
			if !skipDuration {
				m.requestDuration.WithLabelValues(r.Method, path).Observe(duration)
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}
