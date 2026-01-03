package domain

import (
	"ohp/internal/domain/auth"
	"ohp/internal/domain/endpoint"
	"ohp/internal/domain/notifications"
	"ohp/internal/domain/push"
	"ohp/internal/domain/token"
	"ohp/internal/domain/user"

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
