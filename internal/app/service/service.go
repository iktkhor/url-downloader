package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) DownloadFromURLs(URLs []string) error {
    // Отправляем GET-запрос
	wg := sync.WaitGroup{}
	//var err string

	for i, v := range URLs {
		wg.Add(1)
		filename := fmt.Sprintf("image%d.jpeg", i + 1)

		go func() {
			defer wg.Done()
			downloadImage(v, filename)
		}()
	}

    fmt.Println("Waiting")
	wg.Wait()
    fmt.Println("Done")

	return nil
}

func downloadImage(URL string, filename string) {
	resp, err := http.Get(URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}
}

