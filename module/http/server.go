package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/galihrivanto/kotak/config"
	"github.com/galihrivanto/kotak/db"
	"github.com/galihrivanto/kotak/module"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	db     *db.DB
	srv    *echo.Echo
}

func (s *Server) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(ctx)
	address := fmt.Sprintf("%s:%s", s.cfg.HttpServer.Host, s.cfg.HttpServer.Port)

	fmt.Println("Setup API with base API", s.cfg.HttpServer.APIBase)
	s.setupAPI()

	fmt.Println("Setup Static with static URL", s.cfg.HttpServer.StaticURL)
	s.setupStatic()

	go func() {
		if s.cfg.HttpServer.TLS {
			if err := s.srv.StartTLS(address, s.cfg.HttpServer.CertFile, s.cfg.HttpServer.KeyFile); err != nil {
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		} else {
			if err := s.srv.Start(address); err != nil {
				log.Fatalf("Failed to start HTTP server: %v", err)
			}
		}
	}()

	return nil
}

func (s *Server) setupAPI() {
	api := s.srv.Group(s.cfg.HttpServer.APIBase)

	// Account routes
	api.POST("/accounts", s.createAccount)
	api.GET("/accounts/:id/emails", s.getEmails)
	api.GET("/accounts/:id/emails/:email_id", s.getEmail)
}

func (s *Server) setupStatic() {
	s.srv.Static(s.cfg.HttpServer.StaticURL, s.cfg.HttpServer.StaticDir)
	s.srv.Static("/", filepath.Join(s.cfg.HttpServer.StaticDir, "index.html"))
}

func (s *Server) Close() error {
	fmt.Println("Stopping HTTP server")
	s.cancel()
	return s.srv.Shutdown(s.ctx)
}

func NewServer(cfg *config.Config, db *db.DB) *Server {
	svc := &Server{cfg: cfg, db: db}

	svc.srv = echo.New()

	svc.srv.Use(middleware.Logger())
	svc.srv.Use(middleware.Recover())
	svc.srv.Use(middleware.CORS())

	svc.srv.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	return svc
}

func init() {
	module.RegisterModule("http", func(config *config.Config, db *db.DB) module.Module {
		return NewServer(config, db)
	})
}
