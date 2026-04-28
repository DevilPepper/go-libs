package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofrs/uuid/v5"
)

func RequestLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			reqID, err := uuid.NewV4()
			if err == nil {
				requestID = reqID.String()
				r.Header.Set("X-Request-ID", requestID)
				w.Header().Set("X-Request-ID", requestID)
			}
		}
		method := r.Method
		if method == "" {
			method = "GET"
		}
		log.Info("---")
		log.Info(fmt.Sprintf("%s %s", method, r.URL.Path), "request_id", requestID)
		log.Debug("",
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
		rw := newResponseWriter(w)
		start := time.Now()
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		log.Info(rw.status, "duration", duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: -1}
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.status == -1 {
		rw.status = statusCode
		rw.ResponseWriter.WriteHeader(statusCode)
	}
}
