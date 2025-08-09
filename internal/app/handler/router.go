package handler

import "net/http"

func (h *Handler) Routes() http.Handler {
	router := http.NewServeMux()
	
	router.HandleFunc("POST /task", h.PostTaskHandler)
	router.HandleFunc("GET /task/{id}/status", h.StatusTaskHandler)
	router.HandleFunc("POST /task/{id}/load", h.LoadTaskHandler)

	return router
}