package main

import "fmt"

// Struct: tipe data buatan sendiri (kumpulan field)
type Todo struct {
	ID    int
	Title string
	Done  bool
}

// Latihan: struct Karyawan
type Karyawan struct {
	Nama    string
	Jabatan string
	Gaji    int
}

func main() {
	// ==== SLICE: daftar dinamis ====
	fmt.Println("1. Slice (daftar yang bisa bertambah):")
	var hobi []string          // slice kosong
	hobi = append(hobi, "Ngoding")
	hobi = append(hobi, "Main game", "Baca buku") // bisa tambah banyak sekaligus
	fmt.Println("  Hobi:", hobi)
	fmt.Println("  Panjang:", len(hobi), "| Item pertama:", hobi[0])

	// Ambil sebagian slice
	duaPertama := hobi[:2]
	fmt.Println("  2 pertama:", duaPertama)

	// ==== MAP: pasangan key-value ====
	fmt.Println("2. Map (key-value):")
	nilai := map[string]int{
		"Budi": 90,
		"Sari": 85,
	}
	nilai["Dewi"] = 78 // tambah data baru
	fmt.Println("  Nilai:", nilai)
	fmt.Println("  Nilai Budi:", nilai["Budi"])

	// Cek key ada atau tidak
	if v, ada := nilai["Andi"]; ada {
		fmt.Println("  Nilai Andi:", v)
	} else {
		fmt.Println("  Andi tidak ada di map")
	}

	// Hapus key
	delete(nilai, "Sari")
	fmt.Println("  Setelah Sari dihapus:", nilai)

	// ==== LOOP range di map & slice ====
	fmt.Println("3. Loop range di map:")
	for nama, skor := range nilai {
		fmt.Printf("  %s = %d\n", nama, skor)
	}

	// ==== STRUCT ====
	fmt.Println("4. Struct:")
	var t1 Todo          // semua field bernilai nol
	t1.ID = 1
	t1.Title = "Belajar Go"
	t1.Done = false

	t2 := Todo{ID: 2, Title: "Bikin API", Done: true} // deklarasi langsung

	fmt.Println("  Todo 1:", t1)
	fmt.Println("  Todo 2:", t2)
	fmt.Println("  Judul todo 2:", t2.Title)

	// ==== SLICE OF STRUCT (persiapan untuk database in-memory di API) ====
	fmt.Println("5. Slice of struct:")
	todos := []Todo{t1, t2}
	for _, t := range todos {
		status := "belum selesai"
		if t.Done {
			status = "selesai"
		}
		fmt.Printf("  [%d] %s (%s)\n", t.ID, t.Title, status)
	}

	// ==== LATIHAN: slice of struct Karyawan ====
	fmt.Println("6. Daftar karyawan:")
	karyawan := []Karyawan{
		{Nama: "Andi", Jabatan: "Programmer", Gaji: 8_000_000},
		{Nama: "Rina", Jabatan: "Designer", Gaji: 7_500_000},
		{Nama: "Budi", Jabatan: "Manager", Gaji: 12_000_000},
	}
	for _, k := range karyawan {
		fmt.Printf("  %s (%s) - Rp%d\n", k.Nama, k.Jabatan, k.Gaji)
	}
}
