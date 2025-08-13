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
	activeTasks int // Number of current tasks executed in programm
}

func New() *Store {
	return &Store{
		tasks: make([]*Task, 0, maxTasks),
	}
}

// Function for post task, returns index and error
func (s *Store) AddTask() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeTasksLen() == maxTasks {
		return 0, &IndexError{}
	}

	s.tasks = append(s.tasks, NewTask())

	s.incActiveTasksLen()

	return len(s.tasks) - 1, nil
}

// Add url to task by index
func (s *Store) AddTaskURL(taskIndex int, URL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
    task, err := s.getTask(taskIndex)
    if err != nil {
        return err
    }
		
	task.urls = append(task.urls, URL)

	return nil
}

// Get Task URLs copied slice
func (s *Store) GetTaskURLs(taskIndex int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
    task, err := s.getTask(taskIndex)
    if err != nil {
        return nil, err
    }

	res := make([]string, len(task.urls))
	copy(res, task.urls)

	return res, nil
}

func (s *Store) GetTaskStatus(taskIndex int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

    task, err := s.getTask(taskIndex)
    if err != nil {
        return 0, err
    }

	return task.status, nil
}

func (s *Store) SetTaskStatus(taskIndex int, status int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

    task, err := s.getTask(taskIndex)
    if err != nil {
        return err
    }

	task.status = status

	return nil
}

// private function get Task by index, requires correct input
func (s *Store) getTask(taskIndex int) (*Task, error) {
    if !s.isTaskIndexCorrect(taskIndex) {
        return nil, errors.New("task index not found")
    }
    return s.tasks[taskIndex], nil
}

// Function returns number of current tasks
func (s *Store) activeTasksLen() int {
	return s.activeTasks
}

func (s *Store) incActiveTasksLen() {
	s.activeTasks++
}

func (s *Store) DecActiveTasksLen() {
	s.activeTasks--
}

// Check if index is correct
func (s *Store) isTaskIndexCorrect(taskIndex int) bool {
	if taskIndex < 0 || taskIndex >= len(s.tasks) {
        return false
    }

	return true
}

// Check if task stores max number of urls
func (s *Store) IsTaskURLsMax(taskIndex int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := s.tasks[taskIndex]

	return len(task.urls) == maxURLs
}
