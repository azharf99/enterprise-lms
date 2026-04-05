# --- Tahap 1: Builder ---
FROM golang:1.26-alpine AS builder

# Atur direktori kerja di dalam kontainer
WORKDIR /app

# Salin file modul dan unduh dependency
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode proyek
COPY . .

# Compile aplikasi menjadi binary statis bernama 'lms-api'
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o lms-api ./cmd/api/main.go

# --- Tahap 2: Runner ---
FROM alpine:latest

WORKDIR /app

# Tambahkan tzdata agar TimeZone Asia/Jakarta berfungsi akurat untuk ujian CBT
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta

# Salin binary hasil kompilasi dari Tahap 1
COPY --from=builder /app/lms-api .

# Buka port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./lms-api"]