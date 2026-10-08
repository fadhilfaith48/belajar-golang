package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"task-api/model"
	"task-api/store"
)

// Handler: struct yang membungkus store.
// Tiap endpoint diwakili satu method.
type Handler struct {
	store        store.TaskStore
	mu           sync.Mutex
	requestCount int
}

// New: constructor handler. Parameternya INTERFACE,
// jadi kita bisa kirim InMemoryStore atau SQLiteStore.
func New(s store.TaskStore) *Handler {
	return &Handler{store: s}
}

// tambahHitungan: menambah counter dengan aman.
// sync.Mutex mencegah dua request menambah counter bersamaan
// (race condition) — mutex mengunci akses ke data bersama.
func (h *Handler) tambahHitungan() {
	h.mu.Lock()
	h.requestCount++
	h.mu.Unlock()
}

// jumlahRequest: membaca counter dengan aman.
func (h *Handler) jumlahRequest() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.requestCount
}

// writeJSON: helper untuk menulis response JSON + status code
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// Batas pagination
const (
	defaultLimit = 20
	maxLimit     = 100
)

// TaskListResponse: respons GET /tasks.
// Data dibungkus dalam objek (bukan array mentah) supaya bisa ikut
// membawa metadata: jumlah total, halaman sekarang, dll.
type TaskListResponse struct {
	Data       []model.Task `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	Limit      int          `json:"limit"`
	TotalPages int          `json:"total_pages"`
}

// parsePagination: membaca ?page= dan ?limit=.
//
// Parameter kosong memakai default (page=1, limit=20). Parameter yang diisi
// tapi bukan angka atau < 1 DITOLAK dengan 400, bukan diabaikan diam-diam —
// supaya salah ketik langsung ketahuan daripada mengembalikan data yang
// tidak sesuai harapan. limit diklem ke maxLimit supaya klien tidak bisa
// meminta seluruh tabel sekaligus.
func parsePagination(r *http.Request) (page, limit int, err error) {
	page, limit = 1, defaultLimit

	if v := r.URL.Query().Get("page"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 1 {
			return 0, 0, errors.New("page harus angka >= 1")
		}
		page = n
	}

	if v := r.URL.Query().Get("limit"); v != "" {
		n, e := strconv.Atoi(v)
		if e != nil || n < 1 {
			return 0, 0, errors.New("limit harus angka >= 1")
		}
		if n > maxLimit {
			n = maxLimit
		}
		limit = n
	}

	return page, limit, nil
}

// GET /tasks — list dengan filter opsional + pagination
// ?status=&kategori=&search=&page=&limit=
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()

	page, limit, err := parsePagination(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status := r.URL.Query().Get("status")
	kategori := r.URL.Query().Get("kategori")
	// Lowercase sekali di luar loop, bukan per item (lebih hemat)
	search := strings.TrimSpace(r.URL.Query().Get("search"))
	prioritas := r.URL.Query().Get("prioritas")
	sortBy := r.URL.Query().Get("sort_by")
	order := strings.ToLower(r.URL.Query().Get("order"))
	if order != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit
	opts := store.ListOptions{
		Status:    status,
		Kategori:  kategori,
		Search:    search,
		Prioritas: prioritas,
		SortBy:    sortBy,
		Order:     order,
		Limit:     limit,
		Offset:    offset,
	}
	hasil, total, err := h.store.ListFiltered(opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	writeJSON(w, http.StatusOK, TaskListResponse{
		Data:       hasil,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

// POST /tasks — tambah tugas baru
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "format JSON tidak valid", http.StatusBadRequest)
		return
	}
	if err := model.Validasi(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dibuat, err := h.store.Add(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, dibuat)
}

// GET /tasks/{id} — detail satu tugas
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID harus angka", http.StatusBadRequest)
		return
	}
	t, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// PUT /tasks/{id} — ubah seluruh data tugas
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID harus angka", http.StatusBadRequest)
		return
	}
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "format JSON tidak valid", http.StatusBadRequest)
		return
	}
	t.ID = id
	if err := model.Validasi(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.store.Update(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

// DELETE /tasks/{id} — hapus tugas
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID harus angka", http.StatusBadRequest)
		return
	}
	if err := h.store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GET /tasks/stats — statistik tugas per status + total request
func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	stats := h.store.Stats()
	stats["total_request"] = h.jumlahRequest()
	writeJSON(w, http.StatusOK, stats)
}
