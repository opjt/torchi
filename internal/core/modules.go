package core

import (
	"torchi/internal/api"
	"torchi/internal/domain"
	"torchi/internal/infrastructure/db"
	"torchi/internal/pkg"

	"go.uber.org/fx"
)

var Modules = fx.Options(

	pkg.Module,
	api.Module,

	domain.Module,

	//infrastructure
	db.Module,
)
