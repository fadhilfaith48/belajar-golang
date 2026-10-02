package model

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
