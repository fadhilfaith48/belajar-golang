# ===== STAGE 1: build =====
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copy dependency file dulu (biar cache layer saat go mod download)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code lalu build binary
COPY . .
RUN CGO_ENABLED=0 go build -o task-api .

# ===== STAGE 2: runtime =====
# Image hasil akhir hanya berisi binary + minimal sistem
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/task-api .

# Buat folder database lebih dulu (dipakai saat volume di-mount)
RUN mkdir -p /app/data

# Port yang dipakai aplikasi
EXPOSE 8080

CMD ["./task-api"]
