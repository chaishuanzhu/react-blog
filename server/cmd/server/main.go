package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-server/internal/config"
	"blog-server/internal/db"
	apphttp "blog-server/internal/http"
	"blog-server/internal/mail"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server exited", "err", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	var logHandler slog.Handler = slog.NewTextHandler(os.Stdout, nil)
	if cfg.IsProduction() {
		logHandler = slog.NewJSONHandler(os.Stdout, nil)
	}
	slog.SetDefault(slog.New(logHandler))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := db.Migrate(cfg.DatabaseDSN); err != nil {
		return err
	}
	slog.Info("database migrations applied")

	pool, err := db.Open(ctx, cfg.DatabaseDSN)
	if err != nil {
		return err
	}
	defer pool.Close()

	var notifier mail.Notifier = mail.Noop{}
	if cfg.SMTP.Enabled() {
		queue := mail.NewQueue(mail.NewSMTPSender(cfg.SMTP), 100)
		defer queue.Close()
		notifier = queue
		slog.Info("email notifications enabled", "smtp", cfg.SMTP.Host, "notify", cfg.NotifyEmail)
	} else {
		slog.Info("email notifications disabled (SMTP_HOST not set)")
	}

	router, err := apphttp.NewRouter(ctx, cfg, pool, notifier)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", cfg.HTTPAddr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
