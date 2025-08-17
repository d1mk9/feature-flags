package http

import (
	"log"

	"feature-flags/pkg/config"
	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app     *fiber.App
	api     huma.API
	service service.Vars // ← интерфейс
	cfg     *config.Config
}

func NewServer(cfg *config.Config, svc service.Vars) *Server { // ← интерфейс
	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	cfgAPI := huma.DefaultConfig("application/json", "utf-8")
	cfgAPI.Info.Title = "Feature Flags API"
	cfgAPI.Info.Version = "1.0.0"

	api := humafiber.New(app, cfgAPI)

	RegisterRoutes(api, svc)

	return &Server{app: app, api: api, service: svc, cfg: cfg}
}

func (s *Server) Run() {
	addr := ":" + s.cfg.HTTPPort
	log.Printf("🚀 Server is running on %s\n", addr)
	if err := s.app.Listen(addr); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}
}
