package store

import (
	"errors"

	"task-api/model"
)

// ErrNotFound: kesalahan standar saat data tidak ditemukan
var ErrNotFound = errors.New("tugas tidak ditemukan")

// TaskStore: INTERFACE — kontrak "cara menyimpan task".
// Handler cukup bergantung pada interface ini, tidak peduli
// implementasinya (memory sekarang, SQLite nanti).
type TaskStore interface {
	Add(t model.Task) (model.Task, error)
	Get(id int) (model.Task, error)
	List() []model.Task
	Update(t model.Task) (model.Task, error)
	Delete(id int) error
	Stats() map[string]int
}

// InMemoryStore: implementasi TaskStore menggunakan slice.
type InMemoryStore struct {
	tasks  []model.Task
	nextID int
}

// NewInMemory: membuat store baru (function sebagai "constructor")
func NewInMemory() *InMemoryStore {
	return &InMemoryStore{nextID: 1}
}

func (s *InMemoryStore) Add(t model.Task) (model.Task, error) {
	t.ID = s.nextID
	s.nextID++
	if t.Status == "" {
		t.Status = model.StatusTodo
	}
	s.tasks = append(s.tasks, t)
	return t, nil
}

func (s *InMemoryStore) Get(id int) (model.Task, error) {
	for _, t := range s.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return model.Task{}, ErrNotFound
}

func (s *InMemoryStore) List() []model.Task {
	return s.tasks
}

func (s *InMemoryStore) Update(u model.Task) (model.Task, error) {
	for i := range s.tasks {
		if s.tasks[i].ID == u.ID {
			s.tasks[i] = u
			return u, nil
		}
	}
	return model.Task{}, ErrNotFound
}

func (s *InMemoryStore) Delete(id int) error {
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (s *InMemoryStore) Stats() map[string]int {
	stats := map[string]int{
		model.StatusTodo:  0,
		model.StatusDoing: 0,
		model.StatusDone:  0,
	}
	for _, t := range s.tasks {
		stats[t.Status]++
	}
	return stats
}
