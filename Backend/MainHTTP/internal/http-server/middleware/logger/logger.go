package logger

import (
	"fmt"
	log "github.com/go-ozzo/ozzo-log"
	"net/http"
	"os"
	"time"
)

func New(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		resultFunc := func(w http.ResponseWriter, r *http.Request) {
			hostname, _ := os.Hostname()
			logString := fmt.Sprintf("SERVER[%s] --- Request from %s: METHOD - %s || PATH - %s", hostname, r.RemoteAddr, r.Method, r.URL.Path)
			logger.Info(logString)
			t1 := time.Now()
			next.ServeHTTP(w, r)

			logString = fmt.Sprintf("SERVER[%s] --- Request from %s completed in %s", hostname, r.RemoteAddr, time.Since(t1))
			logger.Info(logString)
		}

		return http.HandlerFunc(resultFunc)
	}
}
