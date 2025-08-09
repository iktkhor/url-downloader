package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.s.AddTask(); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusTooManyRequests)
	}
}

func (h *Handler) LoadTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoadTask")
	

}

func (h *Handler) StatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("request for id: " + id))
}
