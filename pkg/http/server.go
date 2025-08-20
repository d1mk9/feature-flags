package http

import (
	"context"
	"time"

	"feature-flags/pkg/config"
	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	api     huma.API
	service service.Vars
	cfg     *config.Config
}

func NewServer(cfg *config.Config, svc service.Vars) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	cfgAPI := huma.DefaultConfig("application/json", "utf-8")
	cfgAPI.Info.Title = "Feature Flags API"
	cfgAPI.Info.Version = "1.0.0"

	api := humafiber.New(app, cfgAPI)
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
	RegisterRoutes(api, svc)

	return &Server{app: app, api: api, service: svc, cfg: cfg}
}

func (s *Server) Run(addr string) error {
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
