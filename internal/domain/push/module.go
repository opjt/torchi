package push

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewWaitMap),
	fx.Provide(NewPushService),
)
