package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iktkhor/url-downloader/internal/app/config"
	"github.com/iktkhor/url-downloader/internal/app/handler"
	"github.com/iktkhor/url-downloader/internal/app/service"
	"github.com/iktkhor/url-downloader/internal/app/store"
)

type App struct {
	s *store.Store
	h *handler.Handler
	svc *service.Service
	cfg *config.Config
}

func New() *App {
	store := store.New()
	config := config.New("config/config.yaml")
	service := service.New()
	handler := handler.New(store, config, service)

	return &App{
		s: store,
		h: handler,
		svc: service,
		cfg: config,
	}
}

func (a *App) Run() error {
	router := a.h.NewRouter()

	addr := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
    server := http.Server{
        Addr:    addr,
        Handler: router,
    }

	fmt.Printf("Server listening on %s in %s mode\n", addr, a.cfg.Env)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}

	return nil
}
