package httpsrvr

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

var (
	options cors.Options = cors.Options{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}
	Middlewares []func(http.Handler) http.Handler
)

func RunHTTPServer(addr string, createHandler func(router chi.Router) http.Handler,
	middlewares []func(http.Handler) http.Handler, allowedAddresses []string, apiPrefix string) {

	updateCorsAllowedAddresses(allowedAddresses)

	apiRouter := chi.NewRouter()

	Middlewares = append(Middlewares, middlewares...)
	setMiddlewares(apiRouter)

	rootRouter := chi.NewRouter()
	if strings.TrimSpace(apiPrefix) == "" {
		apiPrefix = "/"
	}
	rootRouter.Mount(apiPrefix, createHandler(apiRouter))
	logrus.Info("Starting HTTP server")

	_ = http.ListenAndServe(addr, rootRouter)
}

func updateCorsAllowedAddresses(addresses []string) {
	options.AllowedOrigins = append(options.AllowedOrigins, addresses...)
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

	for _, m := range Middlewares {
		router.Use(m)
	}
}

func addCorsMiddleware(router *chi.Mux) {
	if len(options.AllowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(options)
	router.Use(corsMiddleware.Handler)
}
