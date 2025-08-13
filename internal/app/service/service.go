package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
)

var allowedExtensions = []string{".jpeg", ".pdf", ".jpg"}

type DownloadedFile struct {
	Name string
	Data []byte
}

type LoadErrors struct {
	mu sync.Mutex
	errors []error
}

func (e *LoadErrors) AddError(err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.errors = append(e.errors, err)
}

type Service struct {}

func New() *Service {
	return &Service{}
}

func (s *Service) DownloadFromURLs(URLs []string, taskIndex int) ([]DownloadedFile, []error) {
	wg := sync.WaitGroup{}
	e := &LoadErrors{}
	files := make([]DownloadedFile, len(URLs))

	for i, v := range URLs {
		wg.Add(1)
		ext := strings.ToLower(path.Ext(v))
		filename := fmt.Sprintf("task_%d_file_%d%s", taskIndex + 1, i, ext)

		go func() {
			defer wg.Done()

			if !isValidFile(filename) {
				e.AddError(fmt.Errorf("file %s has invalid type\nwrong extension", filename))
				return
			}

			data, err := s.downloadFile(v)
			if err != nil {
	 			e.AddError(err)
				return
			}

			files[i] = DownloadedFile{
				Name: fmt.Sprint(filename),
				Data: data,
			}
		}()
	}

	wg.Wait()

	return files, e.errors
}

func (s *Service) downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func isValidFile(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}