package observability

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	requests uint64
	errors   uint64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (metrics *Metrics) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r)
		atomic.AddUint64(&metrics.requests, 1)
		if recorder.status >= http.StatusInternalServerError {
			atomic.AddUint64(&metrics.errors, 1)
		}
		_ = start
	})
}

func (metrics *Metrics) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = fmt.Fprintf(w, "fusion_http_requests_total %d\n", atomic.LoadUint64(&metrics.requests))
		_, _ = fmt.Fprintf(w, "fusion_http_errors_total %d\n", atomic.LoadUint64(&metrics.errors))
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (recorder *statusRecorder) WriteHeader(status int) {
	recorder.status = status
	recorder.ResponseWriter.WriteHeader(status)
}
