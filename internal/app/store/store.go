package store

import (
	"errors"
	"fmt"
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
	currTasks int // Number of current tasks executed in programm
}

func New() *Store {
	return &Store{
		tasks: make([]*Task, 0, maxTasks),
	}
}

func (s *Store) AddTask() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentStoreLen() == maxTasks {
		return &IndexError{}
	}

	s.tasks = append(s.tasks, NewTask())

	s.currTasks++

	return nil
}

func (s *Store) AddTaskURL(taskIndex int, URL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
    if !s.isTaskIndexCorrect(taskIndex) {
        return errors.New("task index not found")
    }

	task := s.tasks[taskIndex]
	task.urls = append(task.urls, URL)

	return nil
}

func (s *Store) GetTaskURLs(taskIndex int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
    if !s.isTaskIndexCorrect(taskIndex) {
        return nil, errors.New("task index not found")
    }

	task := s.tasks[taskIndex]
	res := make([]string, len(task.urls))
	copy(res, task.urls)

	return res, nil
}

func (s *Store) GetTaskStatus(taskIndex int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

    if !s.isTaskIndexCorrect(taskIndex) {
        return 0, errors.New("task index not found")
    }

	task := s.tasks[taskIndex]

	return task.status, nil
}

func (s *Store) SetTaskStatus(taskIndex int, status int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

    if !s.isTaskIndexCorrect(taskIndex) {
        return errors.New("task index not found")
    }

	task := s.tasks[taskIndex]
	task.status = status

	return nil
}

// Function returns number of current tasks
func (s *Store) currentStoreLen() int {
	return s.currTasks
}

func (s *Store) isTaskIndexCorrect(taskIndex int) bool {
	if taskIndex < 0 || taskIndex >= len(s.tasks) {
        return false
    }

	return true
}

func (s *Store) IsTaskURLsMax(taskIndex int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := s.tasks[taskIndex]

	return len(task.urls) == maxURLs
}
