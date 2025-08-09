package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PostTask")
	h.s.AddTask()
	
}

func (h *Handler) LoadTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoadTask")
}

func (h *Handler) StatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("StatusTask")

	id := r.PathValue("id")
	w.Write([]byte("request for id: " + id))
}

