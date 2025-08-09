package main

import (
	"fmt"
	"net/http"

	"github.com/iktkhor/url-downloader/internal/app/handler"
	"github.com/iktkhor/url-downloader/internal/app/store"
)




func main() {
	store := store.New()
	handler := handler.New(store)
	router := http.NewServeMux()
	
	router.HandleFunc("POST /task", handler.PostTaskHandler)
	router.HandleFunc("GET /task/{id}/status", handler.StatusTaskHandler)
	router.HandleFunc("POST /task/{id}/load", handler.LoadTaskHandler)

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Println("Server listening")
	server.ListenAndServe()
}
