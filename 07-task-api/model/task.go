package model

import (
	"errors"
	"strings"
	"time"
)

// Task: representasi data tugas dalam JSON
type Task struct {
	ID        int    `json:"id"`
	Judul     string `json:"judul"`
	Kategori  string `json:"kategori"`
	Prioritas string `json:"prioritas"`
	Status    string `json:"status"`
	Deadline  string `json:"deadline,omitempty"`
}

// Status yang diperbolehkan
const (
	StatusTodo  = "todo"
	StatusDoing = "doing"
	StatusDone  = "done"
)

// Prioritas yang diperbolehkan
const (
	PrioritasTinggi = "tinggi"
	PrioritasSedang = "sedang"
	PrioritasRendah = "rendah"
)

// MarkDone: method dengan POINTER receiver.
// Mengubah data yang sebenarnya, bukan salinan.
func (t *Task) MarkDone() {
	t.Status = StatusDone
}

// IsDone: method dengan VALUE receiver.
// Hanya membaca, tidak mengubah apapun.
func (t Task) IsDone() bool {
	return t.Status == StatusDone
}

// ValidStatus: cek apakah status diperbolehkan
func ValidStatus(s string) bool {
	return s == StatusTodo || s == StatusDoing || s == StatusDone
}

// ValidPrioritas: cek apakah prioritas diperbolehkan
func ValidPrioritas(p string) bool {
	return p == PrioritasTinggi || p == PrioritasSedang || p == PrioritasRendah
}

// NormalisasiStatus: isi default "todo" bila kosong atau hanya spasi.
func NormalisasiStatus(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return StatusTodo
	}
	return s
}

// NormalisasiPrioritas: isi default "sedang" bila kosong atau hanya spasi.
func NormalisasiPrioritas(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return PrioritasSedang
	}
	return p
}

// Validasi: sekaligus normalisasi lalu cek isi. Dipakai Create dan Update
// agar keduanya berlaku IDENTIK — sebelumnya hanya Create yang memvalidasi.
func Validasi(t *Task) error {
	t.Status = NormalisasiStatus(t.Status)
	t.Prioritas = NormalisasiPrioritas(t.Prioritas)

	if strings.TrimSpace(t.Judul) == "" {
		return errors.New("judul tidak boleh kosong")
	}
	if !ValidStatus(t.Status) {
		return errors.New("status tidak valid")
	}
	if !ValidPrioritas(t.Prioritas) {
		return errors.New("prioritas tidak valid")
	}
	if !ValidDeadline(t.Deadline) {
		return errors.New("deadline tidak valid")
	}
	return nil
}


// ValidDeadline memvalidasi format deadline.
// - Kosong dibolehkan (deadline opsional).
// - Jika diisi, harus dalam format YYYY-MM-DD (mis. 2026-10-08).
func ValidDeadline(s string) bool {
	if strings.TrimSpace(s) == "" {
		return true
	}
	_, err := time.Parse("2006-01-02", s)
	return err == nil
}
