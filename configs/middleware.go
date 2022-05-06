package configs

import (
	"context"
	"net/http"

	"gitlab.com/kaushikayanam/base/context/ctxtypes"
)

func (c *Config) setCloudPlatform(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ctxtypes.CtxVendorKey, c.PlatformConfig.Name)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func setClientDetails(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientId := r.Header.Get("X-Client-Id")
		if clientId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, ctxtypes.CtxClientIdKey, clientId)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (c *Config) GetHttpMiddlewares() (m []func(next http.Handler) http.Handler) {
	for _, p := range c.httpMiddlewareProviders {
		m = append(m, p.GetMiddlewares()...)
	}
	m = append(m, setClientDetails)
	m = append(m, c.setCloudPlatform)
	return
}
