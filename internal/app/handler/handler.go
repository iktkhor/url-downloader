package handler

import (
	"github.com/iktkhor/url-downloader/internal/app/service"
	"github.com/iktkhor/url-downloader/internal/app/store"
)

type Service interface {
	DownloadFromURLs([]string, int) ([]service.DownloadedFile, []error)
}

type Handler struct {
	s *store.Store
	svc Service
}

func New(s *store.Store, svc Service) *Handler {
	return &Handler{
		s: s,
		svc: svc,
	}
}