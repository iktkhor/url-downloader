package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

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

func (s *Service) DownloadFromURLs(URLs []string) ([]DownloadedFile, []error) {
	wg := sync.WaitGroup{}
	e := &LoadErrors{}
	files := make([]DownloadedFile, len(URLs))

	for i, v := range URLs {
		wg.Add(1)
		filename := fmt.Sprintf("image%d.jpeg", i + 1)

		go func() {
			defer wg.Done()
			data, err := downloadFile(v)
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

    fmt.Println("Waiting")
	wg.Wait()
    fmt.Println("Done")

	return files, e.errors
}

func downloadFile(url string) ([]byte, error) {
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
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))
	return buf.Bytes(), nil
}

// createZipFromMemory — архивирует данные из памяти в ZIP
func createZipFromMemory(files []DownloadedFile, zipPath string) error {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for _, file := range files {
		if len(file.Data) == 0 {
			continue // пропускаем пустые (ошибочные) скачивания
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

	if err := zipWriter.Close(); err != nil {
		return err
	}

	return os.WriteFile(zipPath, buf.Bytes(), 0644)
}