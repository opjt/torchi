//go:build standalone

package api

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
	"torchi/internal/pkg/config"
)

//go:embed all:static
var staticFiles embed.FS

func staticHandler() http.Handler {
	sub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return http.NotFoundHandler()
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		if _, err := fs.Stat(sub, path); err != nil {
			r.URL.Path = "/"
		}

		fileServer.ServeHTTP(w, r)
	})
}

func envJsHandler(env config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Write([]byte(`window.__APP_CONFIG__ = {
  VAPID_KEY: "` + env.Vapid.PublicKey + `",
  GITHUB_CLIENT_ID: "` + env.Github.ClientID + `"
};`))
	}
}
