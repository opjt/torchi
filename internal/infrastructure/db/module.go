package db

import (
	"torchi/internal/infrastructure/db/postgresql"

	"go.uber.org/fx"
)

var Module = fx.Options(
	postgresql.Module,
)
