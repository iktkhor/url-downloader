package store

import "fmt"

const maxURLs = 3

type Task struct {
	status int
	urls   []string
}

func (t *Task) AddURL(URL string) {
	fmt.Println("Add URL")

	if len(t.urls) == maxURLs {

	}
}
