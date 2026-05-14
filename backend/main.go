package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

//go:embed all:dist
var dist embed.FS

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	})))
	mux := http.NewServeMux()

	viteProxy := os.Getenv("VITE_DEV_PROXY")
	if viteProxy != "" {
		target, err := url.Parse(viteProxy)
		if err != nil {
			slog.Error("invalid env", "name", "VITE_DEV_PROXY", "value", viteProxy, "error", err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = target.Host
		}

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
		slog.Info("Development mode", "proxy_destination", viteProxy)
	} else {
		staticFS, err := fs.Sub(dist, "dist")
		if err != nil {
			slog.Error("failed to create sub fs", "error", err)
		}
		fileServer := http.FileServer(http.FS(staticFS))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := staticFS.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
				r.URL.Path = "/"
			}
			fileServer.ServeHTTP(w, r)
		})
	}

	addr := "0.0.0.0:8080"
	slog.Info("Server starting", "address", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "error", err)
		return
	}
}
