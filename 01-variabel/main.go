package main

import "fmt"

func main() {
	// 1. Deklarasi variabel dengan tipe eksplisit
	var nama string = "Budi"
	var umur int = 25

	// 2. Deklarasi tanpa tipe (Go menebaknya sendiri)
	var kota = "Jakarta"

	// 3. Deklarasi singkat (paling umum dipakai)
	tinggi := 175.5
	aktif := true

	// 4. Deklarasi banyak variabel sekaligus
	panjang, lebar := 10, 5

	// 5. Konstanta (nilai tidak bisa diubah)
	const phi = 3.14

	fmt.Println("Nama:", nama)
	fmt.Println("Umur:", umur)
	fmt.Println("Kota:", kota)
	fmt.Println("Tinggi:", tinggi, "cm")
	fmt.Println("Aktif:", aktif)
	fmt.Println("Luas persegi panjang:", panjang*lebar)
	fmt.Println("Nilai phi:", phi)
}
