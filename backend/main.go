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
	"time"
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

		proxy := &httputil.ReverseProxy{
			Rewrite: func(pr *httputil.ProxyRequest) {
				pr.SetURL(target)
				pr.In.Host = target.Host
			},
		}

		mux.Handle("/", proxy)
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

	server := &http.Server{
		Handler:           mux,
		Addr:              "0.0.0.0:8080",
		ReadHeaderTimeout: time.Second,
	}

	slog.Info("Server starting", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server error", "error", err)
		return
	}
}
