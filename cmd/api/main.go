// Command api arranca el servidor HTTP de la biblioteca: carga
// configuración, conecta a PostgreSQL, aplica migraciones y sirve la API
// REST hasta recibir una señal de apagado.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/andrew-fe/golang-back/internal/config"
	"github.com/andrew-fe/golang-back/internal/service"
	"github.com/andrew-fe/golang-back/internal/store"
	"github.com/andrew-fe/golang-back/internal/transport/rest"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("el servidor terminó con error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	_ = config.LoadDotEnv(".env")

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := store.NewPostgresDB(cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := store.Migrate(db); err != nil {
		return err
	}

	bookStore := store.NewBookStore(db)
	bookService := service.NewBookService(bookStore)
	handler := rest.NewRouter(bookService, db)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		slog.Info("servidor escuchando", "addr", srv.Addr)
		slog.Info("endpoints disponibles",
			"GET", "/api/v1/books",
			"POST", "/api/v1/books",
			"GET_BY_ID", "/api/v1/books/{id}",
			"PUT", "/api/v1/books/{id}",
			"DELETE", "/api/v1/books/{id}",
			"health", "/healthz",
			"ready", "/readyz",
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrCh:
		return err
	case <-ctx.Done():
		slog.Info("señal de apagado recibida, cerrando servidor...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	slog.Info("servidor apagado correctamente")
	return nil
}
