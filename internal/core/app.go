package core

import (
	"context"
	"log/slog"
	"os"
	"time"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/log"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func RunServer(opt ...fx.Option) {

	env, err := config.NewEnv()
	if err != nil {
		slog.Error("Failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	var logger *log.Logger

	opts := fx.Options(
		fx.Supply(env), // env provide.
		fx.Provide(
			log.NewLogger,
		),
		fx.Options(opt...),
		fx.Populate(&logger), // logger inject
		fx.WithLogger(func(logger *log.Logger) fxevent.Logger {
			return log.NewFxLogger(logger.Logger)
		}),

		fx.Invoke(run),
	)
	app := fx.New(opts)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		logger.Fatal(err)
	}

	// signal wait
	<-app.Done()
	logger.Info("Shutdown signal received. Stopping server...")

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Error("Failed to stop app gracefully",
			slog.String("error", err.Error()),
		)
	}

}

func run(lc fx.Lifecycle, env config.Env, logger *log.Logger) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Service OnStart hook")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Service OnStop hook")

			return nil
		},
	})
}
