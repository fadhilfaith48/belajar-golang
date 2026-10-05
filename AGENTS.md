# Belajar Golang — Progress

Proyek belajar Go dari nol sampai aplikasi full-stack kecil.
Pemakai: pemula total (Windows 11 + WSL Ubuntu).

## Progress saat ini (~85%)

> Rencana fitur lengkap ada di [RENCANA.md](RENCANA.md) (4 fase, dengan checklist).

- [x] 01-variabel — variabel, tipe data
- [x] 02-fungsi — fungsi, parameter, return, multiple return
- [x] 03-loop — for, if/else, switch, break/continue
- [x] 04-struktur-data — slice, map, struct, slice of struct
- [x] 05-error — error sebagai nilai, fmt.Errorf, errors.Is
- [x] 06-api-todo — API To-Do List (net/http, JSON, tanpa dependency) + go test
- [x] 07-task-api — API Manajemen Tugas (project level menengah)
      - SQLite (modernc.org/sqlite), repository pattern (interface TaskStore)
      - Middleware logging (goroutine+channel) & recovery (panic-safe)
      - Request counter (sync.Mutex), table-driven test (12+ test, semua PASS)
      - Dockerfile multi-stage, image task-api:latest (35.7MB)
- [x] 08-task-frontend — Frontend React + Vite (halaman Manajemen Tugas)
      - Terhubung ke backend via proxy Vite (/api -> localhost:8080)
      - Teruji end-to-end: React -> /api -> Go API -> SQLite
- [x] 09-fitur-api — Validasi prioritas, pencarian judul, pagination
      - model: ValidPrioritas, NormalisasiStatus/Prioritas, Validasi() (dipakai Create+Update)
      - BUG FIX: PUT /tasks/{id} sebelumnya TIDAK memvalidasi apa pun
      - GET /tasks?search= (case-insensitive, bisa digabung filter)
      - GET /tasks?page=&limit= (default 1/20, max limit 100, clamp)
      - PENTING: bentuk JSON GET /tasks BERUBAH -> {"data":[], "total", "page", "limit", "total_pages"}
- [x] 10-fase-1-wajah — CORS, PORT dari env, /health, http.Server + timeout, graceful shutdown
      - middleware.CORS() di middleware/middleware.go, tangani preflight OPTIONS (204)
      - main.go: PORT dari env (default 8080) — WAJAB untuk Render
      - main.go: GET /health -> {"status":"ok"} untuk health check Render
      - main.go: http.Server + ReadHeader/Read/Write/Idle timeout (anti slowloris)
      - main.go: graceful shutdown SIGINT/SIGTERM + srv.Shutdown() + tutup store via type assertion
      - PENTING: `var st store.TaskStore = db` — type assertion hanya bisa pada interface

## Git & GitHub

- Repo: https://github.com/fadhilfaith48/belajar-golang (PRIVATE, monorepo 01-08)
- branch `main` sudah tracking `origin/main`
- 40 file ter-commit. Binary hasil `go build` (tanpa ekstensi) sudah masuk .gitignore
- Checkpoint bersih: `261a123` (email `fadhilfaith2@gmail.com`, tanpa trailer Claude)
- Root Directory saat deploy: Render -> `07-task-api`, Vercel -> `08-task-frontend`

## Cara menjalankan

Backend (wsl):
```
wsl
cd "/mnt/d/ALL FOLDER PKL HEXA/Belajar Golang/07-task-api"
go run .
```
Frontend (Windows terminal, folder 08-task-frontend):
```
npm run dev
```
Buka http://localhost:5173. API langsung bisa diakses di http://localhost:8080/tasks.

Docker (opsional, perlu Docker Desktop menyala):
```
docker run -d --name task-api -p 8080:8080 task-api
```

## Lingkungan penting

- Go 1.26.4 terpasang di WSL: `~/.local/go/bin/go` (tanpa sudo; PATH sudah di ~/.bashrc)
- Node.js v24 di Windows (bukan di WSL) — frontend dijalankan dari Windows
- Docker Desktop 28.4 di Windows (engine kadang mati, perlu dinyalakan manual)
- SQLite driver: `modernc.org/sqlite` (murni Go, tanpa CGO/gcc)
  - PENTING: driver ini TIDAK bisa buka file DB di volume Docker Desktop
    (named/anonymous/bind). DB hanya jalan di layer writable container.
- Git Bash lama lag → sekarang semua dijalankan lewat WSL

## Langkah berikutnya (belum dikerjakan)

- [ ] Deploy API ke internet (Render) + Postgres gratis (Neon)
- [ ] Upgrade DB ke PostgreSQL (perlu refactor: pagination di handler -> SQL LIMIT/OFFSET)
- [ ] Login/auth (JWT)
- [ ] Sorting (?sort_by=&order=), validasi format deadline

## Catatan pagination

Filter + search + pagination sekarang dikerjakan di handler (memori Go) dengan
memotong slice. Interface `TaskStore` masih `List() []model.Task`.
Saat pindah ke PostgreSQL, ini perlu diubah ke `List(limit, offset) + Count()`
supaya query memakai `LIMIT`/`OFFSET` di SQL dan tidak menarik seluruh tabel.

## Gaya kerja

- Materi + latihan, bertahap satu langkah per "lanjut"
- Penjelasan kode baris per baris dalam Bahasa Indonesia
- Setiap perubahan diverifikasi (go build, go test, curl)
