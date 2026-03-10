//go:build !standalone

package api

import (
	"net/http"
	"torchi/internal/pkg/config"
)

func staticHandler() http.Handler {
	return nil
}

func envJsHandler(_ config.Env) http.HandlerFunc {
	return nil
}
