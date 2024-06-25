package middleware

import (
	"log"
	"net"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		duration := time.Since(start)
		host, port, _ := net.SplitHostPort(r.RemoteAddr)
		log.Printf(
			"[%s] %s %s:%s %d %.2fms",
			r.Method,
			r.RequestURI,
			host,
			port,
			sw.status,
			float64(duration.Microseconds())/1000.0,
		)
	}

}

type statusWriter struct {
	http.ResponseWriter
	status  int
	written bool
}

func (w *statusWriter) WriterHeader(status int) {
	w.status = status
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}
