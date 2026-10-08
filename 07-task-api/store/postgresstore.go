package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"task-api/model"
)

// PostgresStore implementasi TaskStore untuk PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgres membuat koneksi PostgreSQL baru.
func NewPostgres(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	s := &PostgresStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *PostgresStore) migrate() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS tasks (
	id SERIAL PRIMARY KEY,
	judul TEXT NOT NULL,
	kategori TEXT NOT NULL,
	prioritas TEXT NOT NULL,
	status TEXT NOT NULL,
	deadline TEXT
)`)
	return err
}

func (s *PostgresStore) Add(t model.Task) (model.Task, error) {
	if t.Status == "" {
		t.Status = model.StatusTodo
	}
	if t.Prioritas == "" {
		t.Prioritas = model.PrioritasSedang
	}
	if t.Kategori == "" {
		t.Kategori = "umum"
	}

	var id int
	err := s.db.QueryRow(`
INSERT INTO tasks (judul, kategori, prioritas, status, deadline)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`, t.Judul, t.Kategori, t.Prioritas, t.Status, t.Deadline).Scan(&id)
	if err != nil {
		return model.Task{}, err
	}
	t.ID = id
	return t, nil
}

func (s *PostgresStore) Get(id int) (model.Task, error) {
	var t model.Task
	err := s.db.QueryRow(`
SELECT id, judul, kategori, prioritas, status, deadline
FROM tasks
WHERE id = $1`, id).Scan(&t.ID, &t.Judul, &t.Kategori, &t.Prioritas, &t.Status, &t.Deadline)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Task{}, ErrNotFound
	}
	if err != nil {
		return model.Task{}, err
	}
	return t, nil
}

func (s *PostgresStore) List() []model.Task {
	rows, err := s.db.Query(`
SELECT id, judul, kategori, prioritas, status, deadline
FROM tasks
ORDER BY id ASC`)
	if err != nil {
		return []model.Task{}
	}
	defer rows.Close()
	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Judul, &t.Kategori, &t.Prioritas, &t.Status, &t.Deadline); err != nil {
			continue
		}
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *PostgresStore) ListFiltered(opts ListOptions) ([]model.Task, int, error) {
	where := []string{"1=1"}
	args := []any{}
	idx := 1
	if opts.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", idx))
		args = append(args, opts.Status)
		idx++
	}
	if opts.Prioritas != "" {
		where = append(where, fmt.Sprintf("prioritas = $%d", idx))
		args = append(args, opts.Prioritas)
		idx++
	}
	if opts.Kategori != "" {
		where = append(where, fmt.Sprintf("kategori = $%d", idx))
		args = append(args, opts.Kategori)
		idx++
	}
	if opts.Search != "" {
		where = append(where, fmt.Sprintf("LOWER(judul) LIKE $%d", idx))
		args = append(args, "%"+strings.ToLower(opts.Search)+"%")
		idx++
	}
	whereStr := strings.Join(where, " AND ")

	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE "+whereStr, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "asc"
	if opts.Order == "desc" {
		order = "desc"
	}
	col := "id"
	if c, ok := SortColumns[opts.SortBy]; ok && opts.SortBy != "" {
		col = c
	}
	limit := opts.Limit
	offset := opts.Offset
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}

	query := "SELECT id, judul, kategori, prioritas, status, deadline FROM tasks WHERE " + whereStr + " ORDER BY " + col + " " + order
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		idx++
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tasks := []model.Task{}
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.Judul, &t.Kategori, &t.Prioritas, &t.Status, &t.Deadline); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	return tasks, total, nil
}

func (s *PostgresStore) Update(u model.Task) (model.Task, error) {
	res, err := s.db.Exec(`
UPDATE tasks
SET judul=$2, kategori=$3, prioritas=$4, status=$5, deadline=$6
WHERE id=$1`, u.ID, u.Judul, u.Kategori, u.Prioritas, u.Status, u.Deadline)
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

func (s *PostgresStore) Delete(id int) error {
	res, err := s.db.Exec(`DELETE FROM tasks WHERE id=$1`, id)
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

func (s *PostgresStore) Stats() map[string]int {
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

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
