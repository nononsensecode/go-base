package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func NewBaseHttpLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&BaseHttpLogger{logger})
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*BaseHttpLogEntry)
	return entry.Logger
}

type BaseHttpLogger struct {
	Logger *logrus.Logger
}

func (l *BaseHttpLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &BaseHttpLogEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI

	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Info("Request started")

	return entry
}

type BaseHttpLogEntry struct {
	Logger logrus.FieldLogger
}

func (l *BaseHttpLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status":       status,
		"resp_bytes_length": bytes,
		"resp_elapsed":      elapsed.Round(time.Millisecond / 100).String(),
	})
	l.Logger.Info("Request completed.")
}

func (l *BaseHttpLogEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}
