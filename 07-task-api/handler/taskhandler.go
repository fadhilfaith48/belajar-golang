package handler

import (
	"encoding/json"
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

// GET /tasks — list dengan filter opsional ?status=&kategori=
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	status := r.URL.Query().Get("status")
	kategori := r.URL.Query().Get("kategori")

	hasil := make([]model.Task, 0)
	for _, t := range h.store.List() {
		if status != "" && t.Status != status {
			continue
		}
		if kategori != "" && t.Kategori != kategori {
			continue
		}
		hasil = append(hasil, t)
	}
	writeJSON(w, http.StatusOK, hasil)
}

// POST /tasks — tambah tugas baru
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	h.tambahHitungan()
	var t model.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "format JSON tidak valid", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(t.Judul) == "" {
		http.Error(w, "judul tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if t.Status != "" && !model.ValidStatus(t.Status) {
		http.Error(w, "status tidak valid", http.StatusBadRequest)
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
