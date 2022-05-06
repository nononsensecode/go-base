package common

import (
	"context"
	"net/http"

	"gitlab.com/kaushikayanam/base/context/ctxtypes"
)

func setClientDetails(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientId := r.Header.Get("X-Client-Id")
		if clientId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, ctxtypes.CtxClientIdKey, clientId)
		// How to find vendor ?
		ctx = context.WithValue(ctx, ctxtypes.CtxVendorKey, "")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
