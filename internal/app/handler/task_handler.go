package handler

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	w.Write([]byte("Success! URL added"))

	if h.s.IsTaskURLsMax(id) {
		h.s.SetTaskStatus(id, http.StatusProcessing)
		h.downloadFiles(w, id)
		h.s.SetTaskStatus(id, http.StatusOK)
	}
}

func (h *Handler) downloadFiles(w http.ResponseWriter, taskIndex int) {
	URLs, err := h.s.GetTaskURLs(taskIndex)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	// Обработать ОШИБКИ!
	df, errs := h.svc.DownloadFromURLs(URLs, taskIndex)
	if errs != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	

	err = service.SaveFilesAsZip(df, fmt.Sprintf("task_%d_archive.zip", taskIndex))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// // Создаём ZIP в памяти
	// var buf bytes.Buffer
	// err = h.writeZip(&buf, df)
	// if err != nil {
	// 	http.Error(w, "error writing zip", http.StatusInternalServerError)
	// 	return
	// }

	// // Отдаём архив клиенту
	// w.Header().Set("Content-Type", "application/zip")
	// w.Header().Set("Content-Disposition", `attachment; filename="images.zip"`)
	// w.WriteHeader(http.StatusOK)
	// _, _ = w.Write(buf.Bytes())
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
