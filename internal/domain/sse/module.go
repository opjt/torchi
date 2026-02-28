package sse

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewBroker),
	fx.Invoke(func(lc fx.Lifecycle, broker *Broker) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				broker.Shutdown()
				return nil
			},
		})
	}),
)
