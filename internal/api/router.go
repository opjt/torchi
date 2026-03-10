package api

import (
	"torchi/internal/api/handler"
	middle "torchi/internal/api/middleware"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

func NewRouter(
	subscriptionHandler *handler.SubscriptionHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	endpointHandler *handler.EndpointHandler,
	apiHandler *handler.ApiHandler,
	notiHandler *handler.NotiHandler,
	sseHandler *handler.SSEHandler,
	healthHandler *handler.HealthHandler,

	tokenProvider *token.TokenProvider,
	env config.Env,

	limitMiddleware *middle.RateLimiterManager,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middle.CorsMiddleware(env.FrontUrl))

	// 정적 파일 서빙 (standalone 빌드 시에만 활성화)
	if h := staticHandler(); h != nil {
		r.Get("/env.js", envJsHandler(env))
		r.Handle("/*", h)
	}

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", healthHandler.Check)
		r.With(middle.RateLimitMiddleware(limitMiddleware)).Mount("/v1", apiHandler.Routes())

		r.Mount("/auth", authHandler.Routes())

		r.Group(func(r chi.Router) {
			r.Use(middle.AuthMiddleware(tokenProvider))
			r.Mount("/subscriptions", subscriptionHandler.Routes())
			r.Mount("/users", userHandler.Routes())
			r.Mount("/endpoints", endpointHandler.Routes())
			r.Mount("/notifications", notiHandler.Routes())
			r.Mount("/sse", sseHandler.Routes())
		})
	})

	return r
}

var routeModule = fx.Module("router",
	fx.Provide(
		handler.NewSubscriptionHandler,
		handler.NewAuthHandler,
		handler.NewUserHandler,
		handler.NewEndpointHandler,
		handler.NewNotiHandler,
		handler.NewSSEHandler,
		handler.NewHealthHandler,

		// API
		handler.NewApiHandler,
	),

	fx.Provide(NewRouter),
	fx.Provide(
		middle.NewRateLimiterManager,
	),
)
