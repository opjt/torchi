package pkg

import (
	"time"
	"torchi/internal/pkg/config"
	"torchi/internal/pkg/token"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func(env config.Env) *token.TokenProvider {
			// 직접 아규먼트 값 주입
			// TODO: env에서 넣어주도록 개선
			return token.NewTokenProvider(
				env.JWTSecret, // secret
				"torchi-api",  // issuer

				2*time.Hour,     // accessExpiry 2hour
				23*24*time.Hour, // refreshExpiry 23day
				// 3*time.Second,
				// 20*time.Second,
			)
		},
	),
)
