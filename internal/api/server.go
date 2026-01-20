package api

import (
	"context"
	"net/http"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func NewHTTPServer(lc fx.Lifecycle, router *chi.Mux, env config.Env, log *log.Logger) *http.Server {
	addr := ":" + strconv.Itoa(env.Service.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Debug("Server started on " + addr)
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

var Module = fx.Options(
	routeModule,
	fx.Invoke(NewHTTPServer),
)
