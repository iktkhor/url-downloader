package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/iktkhor/url-downloader/internal/app/service"
)

func (h *Handler) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := h.s.AddTask()
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusTooManyRequests)
		return
	}

	fmt.Fprintf(w, "Created new task\nTask id = %d, Status = %d", id, http.StatusCreated)
}

func (h *Handler) StatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))

	status, err := h.s.GetTaskStatus(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Task id = %d, Status = %d", id, status)
}

func (h *Handler) LoadTaskHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        URL    string `json:"url"`
    }

	id, _ := strconv.Atoi(r.PathValue("id"))

	if status, err := h.s.GetTaskStatus(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
        return
	} else if status != http.StatusCreated {
		http.Error(w, errors.New("number of loaded URLs is maximum").Error(), http.StatusBadRequest)
        return
	}

    defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	if err := h.s.AddTaskURL(int(id), req.URL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
	w.Write([]byte("Success! URL added\n"))

	if h.s.IsTaskURLsMax(id) {
		h.s.SetTaskStatus(id, http.StatusProcessing)

		archiveName := fmt.Sprintf("task_%d_archive.zip", id)
		archivePath := filepath.Join("files", archiveName)

		err := h.downloadFiles(w, id, archivePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		w.Write([]byte(""))
		fmt.Fprintf(w, "Success! Files added to zip\nURL for downloading: /files/%s", archiveName)
		
		fmt.Println(archivePath)

		h.s.SetTaskStatus(id, http.StatusOK)
	}
}

func (h *Handler) downloadFiles(w http.ResponseWriter, taskIndex int, archivePath string) error {
	URLs, _ := h.s.GetTaskURLs(taskIndex)

	df, errs := h.svc.DownloadFromURLs(URLs, taskIndex)
	if errs != nil {
		h.handleErrors(w, errs)
	}

	err := service.SaveFilesAsZip(df, archivePath)
	if err != nil {
		return err
	}
	
	return nil
}

func (h *Handler) handleErrors(w http.ResponseWriter, errs []error) {
	for _, v := range errs {
		http.Error(w, v.Error(), http.StatusInternalServerError)
	}
}
