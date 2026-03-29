package middle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func setup() (*Metrics, *prometheus.Registry) {
	reg := prometheus.NewRegistry()
	m := &Metrics{
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
		}, []string{"method", "path", "status"}),
		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
		}, []string{"method", "path"}),
		requestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
		}),
	}
	reg.MustRegister(m.requestsTotal, m.requestDuration, m.requestsInFlight)
	return m, reg
}

func hasDurationSample(reg *prometheus.Registry, path string) bool {
	mfs, _ := reg.Gather()
	for _, mf := range mfs {
		if mf.GetName() != "http_request_duration_seconds" {
			continue
		}
		for _, m := range mf.GetMetric() {
			for _, l := range m.GetLabel() {
				if l.GetName() == "path" && l.GetValue() == path {
					h := m.GetHistogram()
					return h != nil && h.GetSampleCount() > 0
				}
			}
		}
	}
	return false
}

func hasRequestCount(reg *prometheus.Registry, path string) bool {
	mfs, _ := reg.Gather()
	for _, mf := range mfs {
		if mf.GetName() != "http_requests_total" {
			continue
		}
		for _, m := range mf.GetMetric() {
			for _, l := range m.GetLabel() {
				if l.GetName() == "path" && l.GetValue() == path {
					return m.GetCounter().GetValue() > 0
				}
			}
		}
	}
	return false
}

func request(m *Metrics, pattern, url string) {
	r := chi.NewRouter()
	r.Use(m.Middleware())
	r.Get(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	r.ServeHTTP(w, req)
}

func TestDuration_NormalRoute(t *testing.T) {
	m, reg := setup()
	request(m, "/api/health", "/api/health")

	if !hasDurationSample(reg, "/api/health") {
		t.Error("일반 라우트는 duration이 집계되어야 함")
	}
}

func TestDuration_AskExcluded(t *testing.T) {
	m, reg := setup()
	request(m, "/api/v1/push/{token}/ask", "/api/v1/push/abc123/ask")

	if hasDurationSample(reg, "/api/v1/push/{token}/ask") {
		t.Error("/ask는 duration이 집계되면 안 됨")
	}
	if !hasRequestCount(reg, "/api/v1/push/{token}/ask") {
		t.Error("/ask도 requests_total은 집계되어야 함")
	}
}

func TestDuration_SSEExcluded(t *testing.T) {
	m, reg := setup()
	request(m, "/api/sse/notifications", "/api/sse/notifications")

	if hasDurationSample(reg, "/api/sse/notifications") {
		t.Error("/sse/notifications는 duration이 집계되면 안 됨")
	}
	if !hasRequestCount(reg, "/api/sse/notifications") {
		t.Error("/sse/notifications도 requests_total은 집계되어야 함")
	}
}

