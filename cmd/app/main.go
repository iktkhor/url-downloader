package main

import (
	"fmt"
	"net/http"
)

var tasks [3]*Task

type Task struct {
	maxURL int
	urls [3]string
}

func main() {
	router := http.NewServeMux()
	
	router.HandleFunc("POST /task", PostTaskHandler)
	router.HandleFunc("GET /task/{id}/status", StatusTaskHandler)
	router.HandleFunc("POST /task/{id}/load", LoadTaskHandler)

	server := http.Server{
		Addr: ":8080",
		Handler: router,
	}

	fmt.Println("Server listening")
	server.ListenAndServe()
}


func PostTaskHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("PostTask")
}

func StatusTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("StatusTask")

	id := r.PathValue("id")
	w.Write([]byte("request for id: " + id))
}

func LoadTaskHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoadTask")
}
