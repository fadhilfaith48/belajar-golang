package handler

import (
	"bytes"
	"encoding/json"
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
		{"prioritas tidak valid", `{"judul":"x","prioritas":"ngawur"}`, http.StatusBadRequest},
		{"prioritas kosong diterima", `{"judul":"x","prioritas":""}`, http.StatusCreated},
		{"prioritas tinggi", `{"judul":"x","prioritas":"tinggi"}`, http.StatusCreated},
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
		{"search judul", "?search=belajar", 1},
		{"search huruf besar", "?search=BELAJAR", 1},
		{"search sebagian", "?search=kerja", 2},
		{"search + filter", "?search=kerja&status=done", 1},
		{"search tidak ada", "?search=zzz", 0},
		{"search mengabaikan kategori", "?search=belajar&kategori=kerja", 0},
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

// TestUpdateValidasi: Update sebelumnya TIDAK memvalidasi apa pun,
// sehingga judul bisa dikosongkan dan prioritas bisa berisi ngawur.
func TestUpdateValidasi(t *testing.T) {
	cases := []struct {
		nama   string
		body   string
		status int
	}{
		{"data lengkap", `{"judul":"Update ok","prioritas":"tinggi","status":"doing"}`, http.StatusOK},
		{"judul kosong", `{"judul":""}`, http.StatusBadRequest},
		{"judul spasi saja", `{"judul":"   "}`, http.StatusBadRequest},
		{"status tidak valid", `{"judul":"x","status":"salah"}`, http.StatusBadRequest},
		{"prioritas tidak valid", `{"judul":"x","prioritas":"ngawur"}`, http.StatusBadRequest},
		{"JSON tidak valid", `{bukan json`, http.StatusBadRequest},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			h, st := setup()
			task, err := st.Add(model.Task{Judul: "Awal"})
			if err != nil {
				t.Fatal(err)
			}
			id := strconv.Itoa(task.ID)

			req := httptest.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBufferString(tc.body))
			req.SetPathValue("id", id)
			w := httptest.NewRecorder()
			h.Update(w, req)

			if w.Code != tc.status {
				t.Errorf("kode = %d, ingin %d (body: %s)", w.Code, tc.status, w.Body.String())
			}
		})
	}
}

// TestUpdateMengisiDefault: field kosong harus ternormalisasi,
// bukan tersimpan sebagai string kosong.
func TestUpdateMengisiDefault(t *testing.T) {
	h, st := setup()
	task, err := st.Add(model.Task{Judul: "Awal"})
	if err != nil {
		t.Fatal(err)
	}
	id := strconv.Itoa(task.ID)

	req := httptest.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBufferString(`{"judul":"Diubah"}`))
	req.SetPathValue("id", id)
	w := httptest.NewRecorder()
	h.Update(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("kode = %d, ingin 200", w.Code)
	}
	hasil, err := st.Get(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	if hasil.Prioritas != model.PrioritasSedang {
		t.Errorf("prioritas = %q, ingin %q", hasil.Prioritas, model.PrioritasSedang)
	}
	if hasil.Status != model.StatusTodo {
		t.Errorf("status = %q, ingin %q", hasil.Status, model.StatusTodo)
	}
}

// isiTigaTask: mengisi store dengan 3 task agar pagination punya isi.
func isiTigaTask(t *testing.T, st store.TaskStore) {
	t.Helper()
	for _, judul := range []string{"Satu", "Dua", "Tiga"} {
		if _, err := st.Add(model.Task{Judul: judul}); err != nil {
			t.Fatal(err)
		}
	}
}

// TestPagination: pastikan ?page= & ?limit= memotong hasil dengan benar.
func TestPagination(t *testing.T) {
	cases := []struct {
		nama      string
		query     string
		inginData int
		inginPage int
		inginTot  int
		inginHlm  int
	}{
		{"tanpa param pakai default", "", 3, 1, 3, 1},
		{"halaman 1 limit 2", "?page=1&limit=2", 2, 1, 3, 2},
		{"halaman 2 limit 2", "?page=2&limit=2", 1, 2, 3, 2},
		{"halaman terakhir", "?page=3&limit=2", 0, 3, 3, 2},
		{"limit melebihi total", "?limit=50", 3, 1, 3, 1},
		{"halaman lewat jangkauan", "?page=99", 0, 99, 3, 1},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			h, st := setup()
			isiTigaTask(t, st)

			req := httptest.NewRequest(http.MethodGet, "/tasks"+tc.query, nil)
			w := httptest.NewRecorder()
			h.List(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("kode = %d, ingin 200 (body: %s)", w.Code, w.Body.String())
			}

			var hasil TaskListResponse
			if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
				t.Fatalf("respons tidak bisa di-parse: %v", err)
			}
			if len(hasil.Data) != tc.inginData {
				t.Errorf("jumlah data = %d, ingin %d", len(hasil.Data), tc.inginData)
			}
			if hasil.Page != tc.inginPage {
				t.Errorf("page = %d, ingin %d", hasil.Page, tc.inginPage)
			}
			if hasil.Total != tc.inginTot {
				t.Errorf("total = %d, ingin %d", hasil.Total, tc.inginTot)
			}
			if hasil.TotalPages != tc.inginHlm {
				t.Errorf("total_pages = %d, ingin %d", hasil.TotalPages, tc.inginHlm)
			}
		})
	}
}

// TestPaginationLimitDikelem: limit di atas batas diklem ke maxLimit.
func TestPaginationLimitDikelem(t *testing.T) {
	h, st := setup()
	isiTigaTask(t, st)

	req := httptest.NewRequest(http.MethodGet, "/tasks?limit=9999", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("kode = %d, ingin 200", w.Code)
	}
	var hasil TaskListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatal(err)
	}
	if hasil.Limit != maxLimit {
		t.Errorf("limit = %d, ingin %d (harus diklem)", hasil.Limit, maxLimit)
	}
}

// TestPaginationParameterNgawur: param bukan angka / < 1 harus ditolak 400.
func TestPaginationParameterNgawur(t *testing.T) {
	cases := []string{
		"?page=abc",
		"?page=0",
		"?page=-5",
		"?limit=abc",
		"?limit=0",
		"?limit=-1",
	}

	for _, q := range cases {
		t.Run(q, func(t *testing.T) {
			h, st := setup()
			isiTigaTask(t, st)

			req := httptest.NewRequest(http.MethodGet, "/tasks"+q, nil)
			w := httptest.NewRecorder()
			h.List(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("kode = %d, ingin 400 (body: %s)", w.Code, w.Body.String())
			}
		})
	}
}

// TestPaginationDenganFilter: total harus menghitung hasil filter,
// bukan jumlah seluruh tabel.
func TestPaginationDenganFilter(t *testing.T) {
	h, st := setup()
	isiTigaTask(t, st)
	st.Add(model.Task{Judul: "Empat", Status: model.StatusDone})

	req := httptest.NewRequest(http.MethodGet, "/tasks?status=done&limit=1", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	var hasil TaskListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &hasil); err != nil {
		t.Fatal(err)
	}
	// Hanya 1 task yang status=done, walau ada 4 task di tabel.
	if hasil.Total != 1 {
		t.Errorf("total = %d, ingin 1 (hanya yang ter-filter)", hasil.Total)
	}
	if hasil.TotalPages != 1 {
		t.Errorf("total_pages = %d, ingin 1", hasil.TotalPages)
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
