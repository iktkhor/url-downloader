package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iktkhor/url-downloader/internal/app/handler"
	"github.com/iktkhor/url-downloader/internal/app/service"
	"github.com/iktkhor/url-downloader/internal/app/store"
)

type App struct {
	s *store.Store
	h *handler.Handler
	svc *service.Service
}

func New() *App {
	store := store.New()
	service := service.New()
	handler := handler.New(store, service)

	return &App{
		s: store,
		h: handler,
		svc: service,
	}
}

func (a *App) Run() error {
	router := a.h.Router()

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Println("Server listening")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}

	return nil
}
