package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.s.AddTask(); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusTooManyRequests)
	}
}

func (h *Handler) LoadTaskHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        URL    string `json:"url"`
    }

	id, _ := strconv.Atoi(r.PathValue("id"))

    defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	h.s.AddTaskURL(int(id), req.URL)

	if h.s.IsTaskURLsMax(id) {
		h.s.SetTaskStatus(id, http.StatusProcessing)
		URLs, err := h.s.GetTaskURLs(id)
		if err != nil {
        	http.Error(w, err.Error(), http.StatusBadRequest)
		}
		
		go h.svc.DownloadFromURLs(URLs)

		
	}
}

func (h *Handler) StatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	status, err := h.s.GetTaskStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(strconv.Itoa(status)))
}
