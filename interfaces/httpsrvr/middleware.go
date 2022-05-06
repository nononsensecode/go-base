package httpsrvr

import "net/http"

type MiddlewareProvider interface {
	GetMiddlewares() []func(http.Handler) http.Handler
}
