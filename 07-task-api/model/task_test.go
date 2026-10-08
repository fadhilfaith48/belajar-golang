package model

import "testing"

func TestValidPrioritas(t *testing.T) {
	cases := []struct {
		nama      string
		prioritas string
		ingin     bool
	}{
		{"tinggi", PrioritasTinggi, true},
		{"sedang", PrioritasSedang, true},
		{"rendah", PrioritasRendah, true},
		{"ngawur", "ngawur", false},
		{"huruf besar ditolak", "Tinggi", false},
		{"kosong ditolak", "", false},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			if dapat := ValidPrioritas(tc.prioritas); dapat != tc.ingin {
				t.Errorf("ValidPrioritas(%q) = %v, ingin %v", tc.prioritas, dapat, tc.ingin)
			}
		})
	}
}

func TestNormalisasiPrioritas(t *testing.T) {
	cases := []struct {
		nama  string
		masuk string
		ingin string
	}{
		{"kosong jadi sedang", "", PrioritasSedang},
		{"spasi jadi sedang", "   ", PrioritasSedang},
		{"dipangkas spasi", "  tinggi  ", PrioritasTinggi},
		{"tidak berubah", "rendah", PrioritasRendah},
		{"nilai ngawur tetap", "ngawur", "ngawur"},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			if dapat := NormalisasiPrioritas(tc.masuk); dapat != tc.ingin {
				t.Errorf("NormalisasiPrioritas(%q) = %q, ingin %q", tc.masuk, dapat, tc.ingin)
			}
		})
	}
}

func TestNormalisasiStatus(t *testing.T) {
	cases := []struct {
		nama  string
		masuk string
		ingin string
	}{
		{"kosong jadi todo", "", StatusTodo},
		{"spasi jadi todo", "  ", StatusTodo},
		{"tidak berubah", StatusDoing, StatusDoing},
		{"nilai ngawur tetap", "ngawur", "ngawur"},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			if dapat := NormalisasiStatus(tc.masuk); dapat != tc.ingin {
				t.Errorf("NormalisasiStatus(%q) = %q, ingin %q", tc.masuk, dapat, tc.ingin)
			}
		})
	}
}

// TestValidasi: pastikan Validasi menormalisasi dulu, lalu menolak nilai ngawur.
func TestValidasi(t *testing.T) {
	cases := []struct {
		nama  string
		task  Task
		ingin string // pesan error, kosong = harus lolos
	}{
		{"lengkap", Task{Judul: "A", Prioritas: PrioritasTinggi, Status: StatusDone}, ""},
		{"prioritas kosong diisi default", Task{Judul: "A"}, ""},
		{"judul kosong", Task{}, "judul tidak boleh kosong"},
		{"judul spasi saja", Task{Judul: "   "}, "judul tidak boleh kosong"},
		{"status ngawur", Task{Judul: "A", Status: "ngawur"}, "status tidak valid"},
		{"prioritas ngawur", Task{Judul: "A", Prioritas: "ngawur"}, "prioritas tidak valid"},
	}

	for _, tc := range cases {
		t.Run(tc.nama, func(t *testing.T) {
			task := tc.task
			err := Validasi(&task)

			if tc.ingin == "" {
				if err != nil {
					t.Fatalf("harus lolos, dapat error: %v", err)
				}
				if task.Status != NormalisasiStatus(task.Status) {
					t.Errorf("status belum ternormalisasi: %q", task.Status)
				}
				return
			}
			if err == nil {
				t.Fatalf("harus error %q, tapi lolos", tc.ingin)
			}
			if err.Error() != tc.ingin {
				t.Errorf("pesan = %q, ingin %q", err.Error(), tc.ingin)
			}
		})
	}
}

// TestValidasiMengisiDefault: pastikan field kosong benar-benar terisi,
// bukan sekadar lolos cek.
func TestValidasiMengisiDefault(t *testing.T) {
	var task Task
	task.Judul = "Tanpa prioritas"
	if err := Validasi(&task); err != nil {
		t.Fatal(err)
	}
	if task.Prioritas != PrioritasSedang {
		t.Errorf("prioritas = %q, ingin %q", task.Prioritas, PrioritasSedang)
	}
	if task.Status != StatusTodo {
		t.Errorf("status = %q, ingin %q", task.Status, StatusTodo)
	}
}

func TestValidDeadline(t *testing.T) {
	// kosong boleh
	if !ValidDeadline("") {
		t.Error("deadline kosong seharusnya valid")
	}
	if !ValidDeadline("2026-10-08") {
		t.Error("2026-10-08 seharusnya valid (YYYY-MM-DD)")
	}
	if ValidDeadline("08-10-2026") {
		t.Error("08-10-2026 seharusnya tidak valid (format salah)")
	}
	if ValidDeadline("abc") {
		t.Error("abc seharusnya tidak valid")
	}
	if ValidDeadline("2026-13-01") {
		t.Error("bulan 13 seharusnya tidak valid")
	}
}
