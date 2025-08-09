package store

import (
	"fmt"
	"net/http"
	"sync"
)

const maxTasks = 3

type IndexError struct {
}

func (e *IndexError) Error() string {
	return fmt.Sprintf("Couldn't find free slot for task.\nMax size of tasks is %d", maxTasks)
}

type Store struct {
	mu    sync.Mutex
	tasks []*Task
}

func New() *Store {
	return &Store{
		tasks: make([]*Task, 0, maxTasks),
	}
}

func (s *Store) AddTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// i, err := s.FindIndex()
	// if err != nil {
	// 	return err
	// }
	if len(s.tasks) == maxTasks {
		return &IndexError{}
	}

	s.tasks = append(s.tasks, &Task{
		status: http.StatusCreated,
		urls: make([]string, 0),
	})

	fmt.Println("Add new Task to slice")

	return nil
}

func (s *Store) FindIndex() (int, error) {
	for i, v := range s.tasks {
		if v == nil {
			return i, nil
		}
	}

	return 0, &IndexError{}
}