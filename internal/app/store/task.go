package store

import (
	"net/http"
)

const maxURLs = 3

type Task struct {
	status int
	urls   []string
}

func NewTask() *Task {
	return &Task{
		status: http.StatusCreated,
		urls: make([]string, 0, maxURLs),
	}
}
