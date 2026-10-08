# Roadmap Fitur — Belajar Golang

Rencana 开发 untuk aplikasi Task Manager (`07-task-api` + `08-task-frontend`).

## Status Singkat

| Fase | Paket | Status |
|---|---|---|
| 1 | WAJAH: CORS, PORT, /health, timeout | ✅ Selesai |
| 2 | TAMPILAN: sorting, filter prioritas, deadline | ✅ Selesai |
| 3 | DATABASE: migrasi PostgreSQL | ✅ Selesai |
| 4 | AUTH: login JWT | ⬜ Belum |

---

## Fase 1 — WAJAH (fondasi cloud)

Tanpa fase ini, frontend di Vercel **tidak bisa** konek ke API di Render.

- [x] Middleware `CORS` + tangani preflight `OPTIONS`
- [x] Baca `PORT` dari env, default `8080`
- [x] Endpoint `GET /health`
- [x] `http.Server` dengan `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`
- [x] Graceful shutdown `SIGINT`/`SIGTERM` + `srv.Shutdown()`
- [x] Tutup store saat shutdown (type assertion `Close()`)
- [x] Test `TestCORS`

> **Catatan verifikasi:** 89 test pass, `go vet` bersih, `gofmt` bersih.
> Smoke test: `/health` 200, preflight `OPTIONS` 204 + header CORS,
> `PORT=8099` honoured, fitur lama (search/pagination/validasi 400) utuh.
> `SIGTERM` tidak bisa disimulasikan di Windows, tapi `SIGINT` (Ctrl+C) dan
> `SIGTERM` (dikirim Render) memakai kode yang sama persis.

## Fase 2 — TAMPILAN

- [ ] `store.ListOptions` + `ListFiltered(opts) ([]model.Task, int, error)` di `TaskStore`
- [x] Implementasi di `InMemoryStore` dan `SQLiteStore`
- [x] Whitelist `SortColumns` (anti SQL injection)
- [x] Sorting `?sort_by=&order=asc|desc`
- [x] Filter prioritas `?prioritas=`
- [x] `model.ValidDeadline` (format `YYYY-MM-DD`) masuk ke `Validasi()`
- [x] Frontend: input deadline, kolom deadline, dropdown prioritas, dropdown urutan
- [x] Test sorting, filter, deadline

## Fase 3 — DATABASE (PostgreSQL)

- [x] `go get github.com/lib/pq`
- [x] `store/postgresstore.go` (`SERIAL`, `RETURNING id`, placeholder `$1`)
- [x] Filter/sort/pagination jadi `WHERE ... ORDER BY ... LIMIT/OFFSET` + `COUNT(*)`
- [x] `main.go` mode ganda: ada `DATABASE_URL` → Postgres, tidak ada → SQLite
- [x] `postgresstore_test.go` dengan `t.Skip` bila `DATABASE_URL` kosong

## Fase 4 — AUTH (JWT)

Keputusan: **semua endpoint dilindungi** login.

- [ ] `go get github.com/golang-jwt/jwt/v5` + `golang.org/x/crypto`
- [ ] `model/user.go` (`User{ID, Email, Nama, PasswordHash}`)
- [ ] `store.UserStore` + InMemory/SQLite/Postgres
- [ ] `auth/auth.go` (`Register`/`Login`, bcrypt, token)
- [ ] `auth/middleware.go` (`Authorization: Bearer`)
- [ ] Route `POST /auth/register`, `POST /auth/login`
- [ ] `src/api.js` helper terpusat (base URL + token)
- [ ] Frontend: halaman login/register, `localStorage`, handle 401

---

## Tugas Manual Kamu (di browser)

Aku tidak bisa login ke akun kamu. Total ±15 menit seumur hidup project.

| # | Kapan | Yang dilakukan |
|---|---|---|
| 1 | Sebelum Fase 3 | Daftar [neon.tech](https://neon.tech) → copy `DATABASE_URL` |
| 2 | Saat deploy | Daftar [render.com](https://render.com) → hubungkan repo |
| 3 | Saat deploy | Daftar [vercel.com](https://vercel.com) → hubungkan repo |

---

## Checkpoint & Rollback

Checkpoint saat ini: **`261a123`** (lokal = GitHub, 87 test pass).

Kembalikan seluruh kode ke checkpoint:

```
git reset --hard 261a123
```

Perubahan tidak di-commit tanpa izin. Jadi checkpoint selalu bisa dipulihkan.

## Dependency Baru (semua murni Go, tanpa CGO)

```
github.com/lib/pq              PostgreSQL driver     (Fase 3)
github.com/golang-jwt/jwt/v5   JWT                  (Fase 4)
golang.org/x/crypto/bcrypt     hash password        (Fase 4)
```