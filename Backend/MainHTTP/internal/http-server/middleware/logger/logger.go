package logger

import (
	"fmt"
	log "github.com/go-ozzo/ozzo-log"
	"net/http"
	"time"
)

func New(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		resultFunc := func(w http.ResponseWriter, r *http.Request) {
			logString := fmt.Sprintf("Request from %s: METHOD - %s || PATH - %s", r.RemoteAddr, r.Method, r.URL.Path)
			logger.Info(logString)
			t1 := time.Now()
			next.ServeHTTP(w, r)

			logString = fmt.Sprintf("Request from %s completed in %s", r.RemoteAddr, time.Since(t1))
			logger.Info(logString)
		}

		return http.HandlerFunc(resultFunc)
	}
}
