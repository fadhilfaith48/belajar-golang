package store

import (
	"errors"
	"sort"
	"strings"

	"task-api/model"
)

// ErrNotFound: kesalahan standar saat data tidak ditemukan
var ErrNotFound = errors.New("tugas tidak ditemukan")

// ListOptions untuk filtering, sorting, dan pagination
type ListOptions struct {
	Status    string
	Kategori  string
	Search    string
	Prioritas string
	SortBy    string
	Order     string // "asc" atau "desc"
	Limit     int
	Offset    int
}

// SortColumns: whitelist kolom yang boleh dipakai untuk ORDER BY.
// Tujuannya mencegah SQL injection lewat parameter sort_by.
var SortColumns = map[string]string{
	"id":       "id",
	"judul":    "judul",
	"prioritas": "prioritas",
	"status":    "status",
	"deadline":  "deadline",
}

// TaskStore: INTERFACE — kontrak "cara menyimpan task".
// Handler cukup bergantung pada interface ini, tidak peduli
// implementasinya (memory, SQLite, Postgres nanti).
type TaskStore interface {
	Add(t model.Task) (model.Task, error)
	Get(id int) (model.Task, error)
	List() []model.Task
	ListFiltered(opts ListOptions) ([]model.Task, int, error)
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


func (s *InMemoryStore) ListFiltered(opts ListOptions) ([]model.Task, int, error) {
	// Duplikasi logika filter ke slice (bisa diterapkan ulang di handler test)
	// Implementasi sederhana: filter semua, sort, potong
	res := make([]model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if opts.Status != "" && t.Status != opts.Status {
			continue
		}
		if opts.Prioritas != "" && t.Prioritas != opts.Prioritas {
			continue
		}
		if opts.Kategori != "" && t.Kategori != opts.Kategori {
			continue
		}
		if opts.Search != "" {
			j := strings.ToLower(t.Judul)
			if !strings.Contains(j, strings.ToLower(opts.Search)) {
				continue
			}
		}
		res = append(res, t)
	}
	// sort
	if opts.SortBy != "" {
		col := opts.SortBy
		asc := opts.Order != "desc"
		sort.Slice(res, func(i, k int) bool {
			if col == "id" {
				if asc {
					return res[i].ID < res[k].ID
				}
				return res[i].ID > res[k].ID
			}
			if col == "judul" {
				if asc {
					return strings.ToLower(res[i].Judul) < strings.ToLower(res[k].Judul)
				}
				return strings.ToLower(res[i].Judul) > strings.ToLower(res[k].Judul)
			}
			if col == "prioritas" {
				if asc {
					return res[i].Prioritas < res[k].Prioritas
				}
				return res[i].Prioritas > res[k].Prioritas
			}
			if col == "status" {
				if asc {
					return res[i].Status < res[k].Status
				}
				return res[i].Status > res[k].Status
			}
			if col == "deadline" {
				d1 := res[i].Deadline
				d2 := res[k].Deadline
				if asc {
					return d1 < d2
				}
				return d1 > d2
			}
			return false
		})
	}
	total := len(res)
	l := opts.Limit
	o := opts.Offset
	if l < 0 {
		l = 0
	}
	if o < 0 {
		o = 0
	}
	if l == 0 || o > len(res) {
		return []model.Task{}, total, nil
	}
	end := o + l
	if end > len(res) {
		end = len(res)
	}
	return res[o:end], total, nil
}
