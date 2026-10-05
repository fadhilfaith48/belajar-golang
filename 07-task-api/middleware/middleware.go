package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// logQueue: channel buffered (kapasitas 100) sebagai antrean pesan log.
// Handler (producer) mengirim pesan ke sini, worker (consumer) yang menulisnya.
var logQueue = make(chan string, 100)

// init: membuat SATU goroutine worker yang membaca channel lalu menulis log.
// Karena hanya satu goroutine yang menulis, tidak ada perebutan akses.
func init() {
	go func() {
		for msg := range logQueue {
			log.Print(msg)
		}
	}()
}

// Logging: mencatat tiap request (method, path, durasi) SECARA ASINKRON.
// Request hanya mengirim pesan ke channel lalu selesai — tidak menunggu
// log selesai ditulis. Ini contoh pola producer-consumer.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logQueue <- fmt.Sprintf("[%s] %s (durasi %v)", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS: mengizinkan browser dari origin lain (mis. Vercel) memanggil API ini.
//
// Tanpa header ini, browser memblokir request lintas domain dari frontend.
// Untuk project belajar, "*" sudah cukup. Di produksi sebaiknya whitelist
// origin spesifik supaya hanya domain sendiri yang boleh mengakses.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Preflight: untuk POST/PUT/DELETE lintas domain, browser mengirim
		// OPTIONS dulu untuk menanyakan "boleh tidak?". Jawab langsung
		// tanpa memanggil handler.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Recovery: menangkap panic dari handler mana pun.
// Tanpa ini, panic membuat seluruh server crash (program berhenti).
// Dengan ini, server tetap hidup dan request dibalas 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC ditangkap: %v\n%s", err, debug.Stack())
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
