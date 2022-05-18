package aws

import "net/http"

// For future use
func (a *AWSConfig) GetMiddlewares() []func(http.Handler) http.Handler {
	return a.httpMiddlewares
}
