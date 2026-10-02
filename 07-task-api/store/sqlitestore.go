package store

import (
	"database/sql"
	"errors"

	"task-api/model"

	_ "modernc.org/sqlite"
)

// SQLiteStore: implementasi TaskStore menggunakan database SQLite.
// Menghubungkan ke interface TaskStore yang sama seperti InMemoryStore.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLite: membuka koneksi database + membuat tabel jika belum ada.
func NewSQLite(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Cek koneksi benar-benar bisa dipakai
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Buat tabel jika belum ada (migrasi sederhana)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			judul     TEXT NOT NULL,
			kategori  TEXT,
			prioritas TEXT,
			status    TEXT NOT NULL DEFAULT 'todo',
			deadline  TEXT
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(t model.Task) (model.Task, error) {
	if t.Status == "" {
		t.Status = model.StatusTodo
	}
	res, err := s.db.Exec(
		"INSERT INTO tasks (judul, kategori, prioritas, status, deadline) VALUES (?, ?, ?, ?, ?)",
		t.Judul, t.Kategori, t.Prioritas, t.Status, t.Deadline,
	)
	if err != nil {
		return model.Task{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return model.Task{}, err
	}
	t.ID = int(id)
	return t, nil
}

func (s *SQLiteStore) Get(id int) (model.Task, error) {
	row := s.db.QueryRow(
		"SELECT id, judul, kategori, prioritas, status, deadline FROM tasks WHERE id = ?",
		id,
	)
	return scanTask(row)
}

func (s *SQLiteStore) List() []model.Task {
	rows, err := s.db.Query(
		"SELECT id, judul, kategori, prioritas, status, deadline FROM tasks ORDER BY id",
	)
	if err != nil {
		return []model.Task{}
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		err := rows.Scan(&t.ID, &t.Judul, &t.Kategori, &t.Prioritas, &t.Status, &t.Deadline)
		if err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *SQLiteStore) Update(u model.Task) (model.Task, error) {
	res, err := s.db.Exec(
		"UPDATE tasks SET judul = ?, kategori = ?, prioritas = ?, status = ?, deadline = ? WHERE id = ?",
		u.Judul, u.Kategori, u.Prioritas, u.Status, u.Deadline, u.ID,
	)
	if err != nil {
		return model.Task{}, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}
	if affected == 0 {
		return model.Task{}, ErrNotFound
	}
	return u, nil
}

func (s *SQLiteStore) Delete(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLiteStore) Stats() map[string]int {
	stats := map[string]int{
		model.StatusTodo:  0,
		model.StatusDoing: 0,
		model.StatusDone:  0,
	}
	rows, err := s.db.Query("SELECT status, COUNT(*) FROM tasks GROUP BY status")
	if err != nil {
		return stats
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var jumlah int
		if err := rows.Scan(&status, &jumlah); err != nil {
			continue
		}
		stats[status] = jumlah
	}
	return stats
}

// scanTask: helper untuk mengubah 1 baris SQL menjadi model.Task.
// Mengembalikan ErrNotFound jika tidak ada data.
func scanTask(row *sql.Row) (model.Task, error) {
	var t model.Task
	err := row.Scan(&t.ID, &t.Judul, &t.Kategori, &t.Prioritas, &t.Status, &t.Deadline)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

// Close: menutup koneksi database.
// Penting dipanggil saat program selesai agar file tidak terkunci.
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
