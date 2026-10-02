package store

import (
	"errors"
	"path/filepath"
	"testing"

	"task-api/model"
)

// setupDB: membuat SQLiteStore dengan database di folder temp.
// t.Cleanup memastikan koneksi ditutup setelah test selesai.
func setupDB(t *testing.T) *SQLiteStore {
	t.Helper()
	s, err := NewSQLite(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestAddDanGet(t *testing.T) {
	s := setupDB(t)

	task, err := s.Add(model.Task{Judul: "Belajar test", Kategori: "belajar"})
	if err != nil {
		t.Fatal(err)
	}
	if task.ID == 0 {
		t.Error("ID harus terisi oleh database")
	}
	if task.Status != model.StatusTodo {
		t.Errorf("status default = %q, ingin %q", task.Status, model.StatusTodo)
	}

	got, err := s.Get(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Judul != "Belajar test" {
		t.Errorf("judul = %q, ingin %q", got.Judul, "Belajar test")
	}
}

func TestGetTidakAda(t *testing.T) {
	s := setupDB(t)

	_, err := s.Get(999)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, ingin ErrNotFound", err)
	}
}

func TestUpdate(t *testing.T) {
	s := setupDB(t)

	task, _ := s.Add(model.Task{Judul: "Sebelum", Status: model.StatusTodo})
	task.Judul = "Sesudah"
	task.Status = model.StatusDone

	updated, err := s.Update(task)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Judul != "Sesudah" || updated.Status != model.StatusDone {
		t.Errorf("update tidak tersimpan: %+v", updated)
	}

	// update ID yang tidak ada harus ErrNotFound
	_, err = s.Update(model.Task{ID: 999, Judul: "x"})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("error = %v, ingin ErrNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	s := setupDB(t)

	task, _ := s.Add(model.Task{Judul: "Hapus"})
	if err := s.Delete(task.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(task.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("setelah delete harus ErrNotFound, dapat %v", err)
	}
	if err := s.Delete(task.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("delete berulang harus ErrNotFound, dapat %v", err)
	}
}

func TestStats(t *testing.T) {
	s := setupDB(t)
	s.Add(model.Task{Judul: "A", Status: model.StatusTodo})
	s.Add(model.Task{Judul: "B", Status: model.StatusDone})
	s.Add(model.Task{Judul: "C", Status: model.StatusDone})
	s.Add(model.Task{Judul: "D", Status: model.StatusDoing})

	stats := s.Stats()
	if stats[model.StatusTodo] != 1 {
		t.Errorf("todo = %d, ingin 1", stats[model.StatusTodo])
	}
	if stats[model.StatusDone] != 2 {
		t.Errorf("done = %d, ingin 2", stats[model.StatusDone])
	}
	if stats[model.StatusDoing] != 1 {
		t.Errorf("doing = %d, ingin 1", stats[model.StatusDoing])
	}
}

func TestPersistensi(t *testing.T) {
	// buka store, simpan data, tutup
	path := filepath.Join(t.TempDir(), "persist.db")
	s1, err := NewSQLite(path)
	if err != nil {
		t.Fatal(err)
	}
	s1.Add(model.Task{Judul: "Tersimpan"})
	s1.Close()

	// buka lagi dengan file yang sama, data harus masih ada
	s2, err := NewSQLite(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s2.Close()

	list := s2.List()
	if len(list) != 1 || list[0].Judul != "Tersimpan" {
		t.Errorf("data tidak persisten, isi: %+v", list)
	}
}
