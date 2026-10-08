package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Wiring: buat store, bungkus dengan handler.
	// Mode ganda: jika DATABASE_URL diset, pakai PostgreSQL; kalau tidak, pakai SQLite (lokal).
	var st store.TaskStore
	dsn := os.Getenv("DATABASE_URL")
	if dsn != "" {
		pg, err := store.NewPostgres(dsn)
		if err != nil {
			log.Fatal(err)
		}
		st = pg
		log.Println("store: PostgreSQL aktif")
	} else {
		db, err := store.NewSQLite(dbPath)
		if err != nil {
			log.Fatal(err)
		}
		st = db
		log.Println("store: SQLite aktif")
	}

	// Tutup database saat program berhenti, supaya file SQLite tidak
	// meninggalkan data belum tersimpan. Close() TIDAK ada di interface
	// TaskStore, jadi dicek dulu apakah store-nya mendukung (type assertion).
	defer func() {
		if c, ok := st.(interface{ Close() error }); ok {
			if err := c.Close(); err != nil {
				log.Printf("gagal menutup database: %v", err)
			}
		}
	}()

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

	// Health check: dipakai platform cloud (Render) untuk mengecek
	// apakah service masih hidup. Harus jawaban cepat dan tidak
	// membaca database.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Middleware dibungkus dari luar: Recovery di paling luar supaya
	// juga menangkap panic dari middleware di dalamnya, lalu CORS
	// (header + preflight), lalu Logging, lalu mux.
	// Urutan berpengaruh: request masuk Recovery → CORS → Logging → handler.
	app := middleware.Recovery(middleware.CORS(middleware.Logging(mux)))

	// Port dari environment variable. Platform cloud (Render, Heroku)
	// selalu menyediakan PORT sendiri, jadi TIDAK boleh di-hardcode.
	// Kalau kosong (lokal), pakai default 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// http.Server dengan timeout: tanpa ini, satu koneksi yang lambat
	// bisa menahan resource server selamanya (serangan slowloris).
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app,
		// Batas hanya untuk membaca header request.
		ReadHeaderTimeout: 5 * time.Second,
		// Batas total untuk membaca seluruh request.
		ReadTimeout: 10 * time.Second,
		// Batas total untuk menulis response.
		WriteTimeout: 15 * time.Second,
		// Batas koneksi idle (keep-alive) sebelum ditutup.
		IdleTimeout: 60 * time.Second,
	}

	// Graceful shutdown: platform cloud mengirim SIGTERM sebelum mematikan
	// container. Tangani sinyal itu: berhenti menerima request baru, lalu
	// tunggu request yang sedang berjalan sampai selesai (maks 10 detik).
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop

		log.Println("shutdown diminta, menunggu request yang sedang berjalan...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown dipaksa: %v", err)
			return
		}
		log.Println("server berhenti dengan rapi")
	}()

	log.Printf("Server jalan di http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}
