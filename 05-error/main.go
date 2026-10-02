package main

import (
	"errors"
	"fmt"
	"strconv"
)

// 1. Fungsi yang mengembalikan error
func bagi(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("tidak bisa membagi dengan nol")
	}
	return a / b, nil
}

// 2. Konversi string ke angka (mengembalikan error)
func konversiAngka(teks string) (int, error) {
	angka, err := strconv.Atoi(teks)
	if err != nil {
		return 0, fmt.Errorf("gagal konversi '%s': %w", teks, err)
	}
	return angka, nil
}

// Latihan: fungsi validasi umur
func cekUmur(umur int) error {
	if umur < 17 {
		return errors.New("belum cukup umur")
	}
	return nil
}

func main() {
	// ==== 1. Menangani error dari fungsi ====
	fmt.Println("1. Error handling:")
	hasil, err := bagi(10, 2)
	if err != nil {
		fmt.Println("  Gagal:", err)
	} else {
		fmt.Println("  10 / 2 =", hasil)
	}

	hasil, err = bagi(10, 0)
	if err != nil {
		fmt.Println("  Gagal:", err)
	} else {
		fmt.Println("  10 / 0 =", hasil)
	}

	// ==== 2. Error dari konversi ====
	fmt.Println("2. Konversi string ke angka:")
	umur, err := konversiAngka("25")
	if err != nil {
		fmt.Println("  Gagal:", err)
	} else {
		fmt.Println("  Umur:", umur, "(tipe", fmt.Sprintf("%T", umur) + ")")
	}

	umur, err = konversiAngka("bukan-angka")
	if err != nil {
		fmt.Println("  Gagal:", err)
	} else {
		fmt.Println("  Umur:", umur)
	}

	// ==== 3. errors.Is untuk cek error spesifik ====
	fmt.Println("3. Cek error spesifik dengan errors.Is:")
	var errBagiNol = errors.New("tidak bisa membagi dengan nol")
	if errors.Is(err, errBagiNol) {
		fmt.Println("  Ini error pembagian nol")
	} else {
		fmt.Println("  Bukan error pembagian nol")
	}

	// ==== 4. LATIHAN: cek umur ====
	fmt.Println("4. Validasi umur:")
	for _, u := range []int{20, 15} {
		if err := cekUmur(u); err != nil {
			fmt.Printf("  Umur %d: ditolak (%v)\n", u, err)
		} else {
			fmt.Printf("  Umur %d: boleh\n", u)
		}
	}
}
