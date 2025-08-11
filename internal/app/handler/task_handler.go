package handler

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/iktkhor/url-downloader/internal/app/service"
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
		go func() {
			h.downloadFiles(w, id)
		}()
	}
}

func (h *Handler) downloadFiles(w http.ResponseWriter, taskIndex int) {
	h.s.SetTaskStatus(taskIndex, http.StatusProcessing)
	URLs, err := h.s.GetTaskURLs(taskIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	data, e := h.svc.DownloadFromURLs(URLs)
	if e != nil {
		fmt.Println(e)
	}

	// Создаём ZIP в памяти
	var buf bytes.Buffer
	err = h.writeZip(&buf, data)
	if err != nil {
		http.Error(w, "Ошибка создания ZIP", http.StatusInternalServerError)
		return
	}

	// Отдаём архив клиенту
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="images.zip"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())

	h.s.SetTaskStatus(taskIndex, http.StatusOK)
}

func (h *Handler) writeZip(w io.Writer, files []service.DownloadedFile) error {
	zipWriter := zip.NewWriter(w)

	for _, file := range files {
		if len(file.Data) == 0 {
			continue
		}
		f, err := zipWriter.Create(file.Name)
		if err != nil {
			return err
		}
		_, err = f.Write(file.Data)
		if err != nil {
			return err
		}
	}

	return zipWriter.Close()
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
