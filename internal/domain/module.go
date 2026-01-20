package domain

import (
	"torchi/internal/domain/auth"
	"torchi/internal/domain/endpoint"
	"torchi/internal/domain/notifications"
	"torchi/internal/domain/push"
	"torchi/internal/domain/token"
	"torchi/internal/domain/user"

	"go.uber.org/fx"
)

var Module = fx.Options(

	push.Module,
	auth.Module,
	user.Module,
	token.Module,
	endpoint.Module,
	notifications.Module,
)
