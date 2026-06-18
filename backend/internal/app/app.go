package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/herrnan/laga/internal/app/web"
	"github.com/herrnan/laga/internal/util"
	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, log *slog.Logger, cfg Config) error {
	eg, ctx := errgroup.WithContext(ctx)

	mux := http.NewServeMux()

	if cfg.ExternalWebServerAddress != "" {
		serveWeb, err := web.ServeProxy(cfg.ExternalWebServerAddress)
		if err != nil {
			return fmt.Errorf("create web server proxy: %w", err)
		}

		log.Info("Running against external web server", "external_web_server_address", cfg.ExternalWebServerAddress)
		mux.HandleFunc("/", serveWeb)
	} else {
		serveWeb, err := web.ServeStatic()
		if err != nil {
			return fmt.Errorf("create static web server: %w", err)
		}

		mux.HandleFunc("/", serveWeb)
	}

	server := &http.Server{
		Handler:           http.MaxBytesHandler(mux, 1<<20 /* 1mb */),
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 1 * time.Second,
	}

	eg.Go(func() error {
		log.Info("Server starting", "address", server.Addr)

		err := util.ListenAndServe(ctx, server, time.Second*10)
		if err != nil {
			return nil
		}

		return fmt.Errorf("http server failed: %w", err)
	})

	return eg.Wait()
}
