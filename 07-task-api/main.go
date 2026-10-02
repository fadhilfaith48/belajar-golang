package main

import (
	"log"
	"net/http"
	"os"

	"task-api/handler"
	"task-api/middleware"
	"task-api/model"
	"task-api/store"
)

func main() {
	// Lokasi database bisa diatur lewat environment variable DB_PATH.
	// Default "tasks.db". Contoh: DB_PATH=/app/data/tasks.db ./task-api
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tasks.db"
	}

	// Wiring: buat store, bungkus dengan handler
	st, err := store.NewSQLite(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	h := handler.New(st)

	// Seed: isi data contoh hanya jika database masih kosong
	if len(st.List()) == 0 {
		for _, judul := range []string{"Belajar variabel", "Belajar fungsi", "Belajar loop"} {
			if _, err := st.Add(model.Task{Judul: judul, Kategori: "belajar", Prioritas: "sedang"}); err != nil {
				log.Fatal(err)
			}
		}
		log.Println("Seed data awal berhasil ditambahkan")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", h.List)
	mux.HandleFunc("POST /tasks", h.Create)
	mux.HandleFunc("GET /tasks/{id}", h.Get)
	mux.HandleFunc("PUT /tasks/{id}", h.Update)
	mux.HandleFunc("DELETE /tasks/{id}", h.Delete)
	mux.HandleFunc("GET /tasks/stats", h.Stats)

	// Middleware dibungkus dari luar: Recovery di paling luar,
	// lalu Logging, lalu mux di tengah.
	// Urutan berpengaruh: request masuk Recovery → Logging → handler.
	app := middleware.Recovery(middleware.Logging(mux))

	log.Println("Server jalan di http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", app))
}
