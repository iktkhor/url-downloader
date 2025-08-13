package handler

import "net/http"

func (h *Handler) NewRouter() http.Handler {
	router := http.NewServeMux()
	
	router.HandleFunc("POST /task", h.PostTaskHandler)
	router.HandleFunc("GET /task/{id}/status", h.StatusTaskHandler)
	router.HandleFunc("POST /task/{id}/load", h.LoadTaskHandler)

	fs := http.FileServer(http.Dir("files"))
	router.Handle("/files/", http.StripPrefix("/files/", fs))

	return router
}