package core

import (
	"ohp/internal/api"
	"ohp/internal/domain"
	"ohp/internal/infrastructure/db"
	"ohp/internal/pkg"

	"go.uber.org/fx"
)

var Modules = fx.Options(

	pkg.Module,
	api.Module,

	domain.Module,

	//infrastructure
	db.Module,
)
