package main

import "fmt"

// 1. Fungsi tanpa parameter dan tanpa return
func sapa() {
	fmt.Println("Halo, selamat belajar Golang!")
}

// 2. Fungsi dengan parameter (return kosong, hanya aksi)
func sapaNama(nama string) {
	fmt.Println("Halo", nama)
}

// 3. Fungsi dengan parameter dan 1 return value
func tambah(a int, b int) int {
	return a + b
}

// 4. Fungsi dengan tipe parameter yang sama (bisa disingkat)
func kurang(a, b int) int {
	return a - b
}

// 5. Fungsi dengan multiple return values
func bagi(a, b int) (int, int) {
	hasil := a / b
	sisa := a % b
	return hasil, sisa
}

// 6. Named return values (nama variabel return dideklarasikan di signature)
func kali(a, b int) (hasil int) {
	hasil = a * b
	return
}

func main() {
	sapa()
	sapaNama("Budi")

	fmt.Println("5 + 3 =", tambah(5, 3))
	fmt.Println("5 - 3 =", kurang(5, 3))

	hasil, sisa := bagi(7, 2)
	fmt.Println("7 / 2 =", hasil, "sisa", sisa)

	fmt.Println("4 * 3 =", kali(4, 3))
}
