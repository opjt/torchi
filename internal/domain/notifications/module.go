package notifications

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewNotiService,
		NewNotiRepository,
	),
)
