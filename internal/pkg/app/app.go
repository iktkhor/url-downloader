package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iktkhor/url-downloader/internal/app/handler"
	"github.com/iktkhor/url-downloader/internal/app/store"
)

type App struct {
	s *store.Store
	h *handler.Handler
}

func New() *App {
	store := store.New()
	handler := handler.New(store)

	return &App{
		s: store,
		h: handler,
	}
}

func (a *App) Run() error {
	router := a.h.Routes()

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal()
	}

	fmt.Println("Server listening")

	return nil
}