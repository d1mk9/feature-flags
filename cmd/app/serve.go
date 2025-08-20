package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"feature-flags/pkg/config"
	httpapi "feature-flags/pkg/http"
	"feature-flags/pkg/repository"
	"feature-flags/pkg/service"
	"feature-flags/pkg/storage"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Запустить HTTP сервер",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoad()

		dbs := storage.MustInitPostgres(cfg.PostgresDSN())
		defer func() {
			if err := dbs.SQL.Close(); err != nil {
				log.Printf("db close error: %v", err)
			}
		}()

		repo := repository.NewPostgresRepository(dbs.Reform)
		svc, err := service.NewFeatureService(repo, 256, 15)
		if err != nil {
			log.Printf("failed to init service: %v", err)
			return
		}

		srv := httpapi.NewServer(cfg, svc)
		addr := ":8080"

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		go func() {
			log.Printf("server is running on %s", addr)
			if err := srv.Run(addr); err != nil {
				errCh <- err
				return
			}
			errCh <- nil
		}()

		select {
		case <-ctx.Done():
			log.Println("shutdown signal received")
		case err := <-errCh:
			if err != nil {
				log.Printf("server run error: %v", err)
			}
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown error: %v", err)
		}

		log.Println("server stopped gracefully")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
