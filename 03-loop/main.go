package main

import "fmt"

func main() {
	// ==== FOR: bentuk 1 — seperti for klasik ====
	fmt.Println("1. Hitung 1 sampai 5:")
	for i := 1; i <= 5; i++ {
		fmt.Println("  i =", i)
	}

	// ==== FOR: bentuk 2 — seperti while ====
	fmt.Println("2. While-style (hitung mundur):")
	n := 3
	for n > 0 {
		fmt.Println("  n =", n)
		n--
	}

	// ==== FOR: bentuk 3 — range (untuk slice/array) ====
	fmt.Println("3. Loop pakai range:")
	nama := []string{"Budi", "Sari", "Dewi"}
	for i, orang := range nama {
		fmt.Println("  Index", i, "=", orang)
	}

	// ==== IF / ELSE ====
	fmt.Println("4. Kondisi if/else:")
	umur := 17
	if umur >= 17 {
		fmt.Println("  Boleh buat KTP")
	} else {
		fmt.Println("  Belum boleh buat KTP")
	}

	// ==== IF dengan statement pendek ====
	fmt.Println("5. If dengan inisialisasi:")
	if skor := 80; skor >= 75 {
		fmt.Println("  Lulus (skor", skor, ")")
	} else {
		fmt.Println("  Tidak lulus")
	}

	// ==== SWITCH ====
	fmt.Println("6. Switch:")
	hari := "Sabtu"
	switch hari {
	case "Senin", "Selasa", "Rabu", "Kamis", "Jumat":
		fmt.Println("  Hari kerja")
	case "Sabtu", "Minggu":
		fmt.Println("  Hari libur")
	default:
		fmt.Println("  Hari tidak dikenal")
	}

	// ==== BREAK & CONTINUE ====
	fmt.Println("7. Break dan continue (bilangan genap dari 1-10, stop di 8):")
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			continue // skip angka ganjil
		}
		if i > 8 {
			break // berhenti setelah 8
		}
		fmt.Println("  ", i)
	}

	// ==== LATIHAN: tabel perkalian 1-5 (2 loop bersarang) ====
	fmt.Println("8. Tabel perkalian 1-5:")
	for a := 1; a <= 5; a++ {
		for b := 1; b <= 5; b++ {
			fmt.Printf("%d x %d = %d\t", a, b, a*b)
		}
		fmt.Println()
	}
}
