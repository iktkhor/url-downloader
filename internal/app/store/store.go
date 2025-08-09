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
	return fmt.Sprintf("Couldn't find free slot. Max size of tasks is %d", maxTasks)
}

type Store struct {
	mu    sync.Mutex
	tasks []*Task
}

func New() *Store {
	return &Store{
		tasks: make([]*Task, maxTasks),
	}
}

func (s *Store) AddTask(t Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	i, err := s.FindIndex()
	if err != nil {
		return err
	}

	s.tasks[i] = &Task{
		status: http.StatusCreated,
	}

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