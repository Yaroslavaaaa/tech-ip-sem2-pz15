package storage

import (
	"errors"
	"sync"
	"time"

	"tasks-app/internal/models"
)

type MemoryStorage struct {
	mu     sync.Mutex
	tasks  []models.Task
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make([]models.Task, 0),
		nextID: 1,
	}
}

func (s *MemoryStorage) GetAll() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.tasks
}

func (s *MemoryStorage) GetByID(id int) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.tasks {
		if task.ID == id {
			t := task
			return &t, nil
		}
	}

	return nil, errors.New("task not found")
}

func (s *MemoryStorage) Create(title, description string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	task := models.Task{
		ID:          s.nextID,
		Title:       title,
		Description: description,
		Done:        false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.tasks = append(s.tasks, task)
	s.nextID++

	return task
}

func (s *MemoryStorage) Update(id int, title, description string, done bool) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks[i].Title = title
			s.tasks[i].Description = description
			s.tasks[i].Done = done
			s.tasks[i].UpdatedAt = time.Now()

			return &s.tasks[i], nil
		}
	}

	return nil, errors.New("task not found")
}

func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}

	return errors.New("task not found")
}
