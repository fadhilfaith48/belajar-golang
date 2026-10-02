package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"task-api/model"
	"task-api/store"
)

// setup: membuat store + handler baru untuk tiap test.
// Mengembalikan store juga supaya test bisa menyiapkan data awal.
func setup() (*Handler, store.TaskStore) {
	st := store.NewInMemory()
	return New(st), st
}

// Table-driven test: satu fungsi test, banyak kasus dalam tabel.
func TestCreate(t *testing.T) {
	h, _ := setup()

	cases := []struct {
		nama   string
		body   string
		status int
	}{
		{"judul valid", `{"judul":"Belajar test","kategori":"belajar"}`, http.StatusCreated},
		{"judul kosong", `{"judul":""}`, http.StatusBadRequest},
		{"JSON tidak valid", `{bukan json`, http.StatusBadRequest},
		{"status tidak valid", `{"judul":"x","status":"salah"}`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(tc.body))
			w := httptest.NewRecorder()
			h.Create(w, req)

			if w.Code != tc.status {
				t.Errorf("kode = %d, ingin %d (body: %s)", w.Code, tc.status, w.Body.String())
			}
		})
	}
}

func TestGet(t *testing.T) {
	h, st := setup()
	task, err := st.Add(model.Task{Judul: "Ada"})
	if err != nil {
		t.Fatal(err)
	}
	existingID := strconv.Itoa(task.ID)

	cases := []struct {
		nama   string
		id     string
		status int
	}{
		{"id ada", existingID, http.StatusOK},
		{"id tidak ada", "999", http.StatusNotFound},
		{"id bukan angka", "abc", http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks/"+tc.id, nil)
			req.SetPathValue("id", tc.id)
			w := httptest.NewRecorder()
			h.Get(w, req)

			if w.Code != tc.status {
				t.Errorf("kode = %d, ingin %d", w.Code, tc.status)
			}
		})
	}
}

func TestListFilter(t *testing.T) {
	h, st := setup()
	st.Add(model.Task{Judul: "Kerjaan done", Status: model.StatusDone, Kategori: "kerja"})
	st.Add(model.Task{Judul: "Kerjaan todo", Status: model.StatusTodo, Kategori: "kerja"})
	st.Add(model.Task{Judul: "Belajar done", Status: model.StatusDone, Kategori: "belajar"})

	cases := []struct {
		nama   string
		query  string
		jumlah int
	}{
		{"tanpa filter", "", 3},
		{"filter status done", "?status=done", 2},
		{"filter kategori kerja", "?kategori=kerja", 2},
		{"filter kombinasi", "?status=done&kategori=belajar", 1},
		{"tidak ada yang cocok", "?status=doing", 0},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tasks"+tc.query, nil)
			w := httptest.NewRecorder()
			h.List(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("kode = %d, ingin 200", w.Code)
			}
			// hitung jumlah objek JSON dengan menghitung kemunculan "id":
			jumlah := strings.Count(w.Body.String(), `"id":`)
			if jumlah != tc.jumlah {
				t.Errorf("jumlah hasil = %d, ingin %d", jumlah, tc.jumlah)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	h, st := setup()
	task, err := st.Add(model.Task{Judul: "Hapus saya"})
	if err != nil {
		t.Fatal(err)
	}
	existingID := strconv.Itoa(task.ID)

	cases := []struct {
		nama   string
		id     string
		status int
	}{
		{"hapus yang ada", existingID, http.StatusNoContent},
		{"hapus yang tidak ada", "999", http.StatusNotFound},
		{"id bukan angka", "abc", http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tc.id, nil)
			req.SetPathValue("id", tc.id)
			w := httptest.NewRecorder()
			h.Delete(w, req)

			if w.Code != tc.status {
				t.Errorf("kode = %d, ingin %d", w.Code, tc.status)
			}
		})
	}
}

func TestStatsMenghitungRequest(t *testing.T) {
	h, st := setup()
	st.Add(model.Task{Judul: "Satu", Status: model.StatusTodo})
	st.Add(model.Task{Judul: "Dua", Status: model.StatusDone})

	// kirim beberapa request agar counter bertambah
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
		w := httptest.NewRecorder()
		h.List(w, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/stats", nil)
	w := httptest.NewRecorder()
	h.Stats(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("kode = %d, ingin 200", w.Code)
	}
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte(`"total_request":4`)) {
		t.Errorf("total_request seharusnya 4, dapat: %s", body)
	}
	if !bytes.Contains([]byte(body), []byte(`"done":1`)) {
		t.Errorf("stats done seharusnya 1, dapat: %s", body)
	}
}
