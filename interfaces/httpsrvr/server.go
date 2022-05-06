package httpsrvr

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

var (
	options cors.Options
)

func RunHTTPServer(addr string, createHandler func(router chi.Router) http.Handler, opts cors.Options, apiPrefix string) {
	options = opts
	apiRouter := chi.NewRouter()
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	rootRouter.Mount(apiPrefix, createHandler(apiRouter))
	logrus.Info("Starting HTTP server")

	_ = http.ListenAndServe(addr, rootRouter)
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	addCorsMiddleware(router)

	router.Use(
		middleware.SetHeader("X-Content-Type-Options", "nosnif"),
		middleware.SetHeader("X-Frame-Options", "deny"),
	)
	router.Use(middleware.NoCache)
}

func addCorsMiddleware(router *chi.Mux) {
	if len(options.AllowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(options)
	router.Use(corsMiddleware.Handler)
}
