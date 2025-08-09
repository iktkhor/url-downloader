package store

import "fmt"

type Task struct {
	status int
	urls   [3]string
}

func (t *Task) AddURL(URL string) {
	fmt.Println("Add URL")
}