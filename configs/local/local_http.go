package local

import "net/http"

// For future use
func (l *LocalConfig) GetMiddlewares() []func(http.Handler) http.Handler {
	return l.httpMiddlewares
}
