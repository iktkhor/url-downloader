package app

import "net/http"

var tasks [3]*Task

type Task struct {
	maxURL int
	urls [3]string
}

func main() {
	
	http.HandleFunc("/task", createTaskHandler)
	http.HandleFunc("/task/{id}/status", statusTaskHandler)
	http.HandleFunc("/task/{id}/load", loadTaskHandler)

}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	default:
	
	}
}

func postTask(w http.ResponseWriter, r *http.Request) error {
	
}

func statusTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func loadTaskHandler(w http.ResponseWriter, r *http.Request) {

}

