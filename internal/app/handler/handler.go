package handler

import (
	"github.com/iktkhor/url-downloader/internal/app/store"
)

type Handler struct {
	s *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{
		s: s,
	}
}