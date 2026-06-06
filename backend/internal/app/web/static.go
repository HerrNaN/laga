package web

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var dist embed.FS

func ServeStatic() (http.HandlerFunc, error) {
	staticFS, err := fs.Sub(dist, "dist")
	if err != nil {
		return nil, fmt.Errorf("create sub fs: %w", err)
	}

	fileServer := http.FileServer(http.FS(staticFS))

	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := staticFS.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r.URL.Path = "/"
		}

		w.Header().Set("Cache-Control", cacheControl(r.URL.Path))

		fileServer.ServeHTTP(w, r)
	}, nil
}

func cacheControl(path string) string {
	switch {
	case strings.HasPrefix(path, "/assets/"):
		return "public, max-age=31536000, immutable"
	case path == "/index.html" || path == "/":
		return "no-cache"
	default:
		return "no-cache"
	}
}
