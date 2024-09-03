package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to run server", "error", err)
	}
}

func run(ctx context.Context) error {
	app, err := setupApplication()
	if err != nil {
		return fmt.Errorf("failed to setup application: %w", err)
	}

	mux, err := app.SetupRoutes()
	if err != nil {
		return fmt.Errorf("failed to setup routes: %w", err)
	}
	slog.InfoContext(ctx, "starting server on :8080")
	return http.ListenAndServe(":8080", mux)
}
