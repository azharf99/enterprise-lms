package utils

import (
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64

	// 1. Cek apakah sudah ada pengguna dengan Email "admin@gmail.com" di database
	db.Model(&domain.User{}).Where("email = ?", "admin@gmail.com").Count(&count)

	// 2. Jika belum ada Admin sama sekali, buat satu!
	if count == 0 {
		fmt.Println("[Seeder] Mendeteksi database baru. Menyiapkan akun Super Admin...")

		// Hash password default (misal: "admin123")
		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			fmt.Println("❌ [Seeder] Gagal hash password Super Admin:", err)
		}

		// Buat model Super Admin
		admin := domain.User{
			Name:     "Admin Ganteng",
			Email:    "admin@gmail.com",
			Password: string(hashedPassword),
			Role:     "Admin",
		}

		// Simpan ke database
		if err := db.Create(&admin).Error; err != nil {
			fmt.Println("❌ [Seeder] Gagal membuat akun Super Admin:", err)
		}

		fmt.Println("[Seeder] ✅ Berhasil membuat Super Admin! (Email: admin@gmail.com | Pass: admin123)")
	} else {
		// Jika sudah ada, abaikan saja agar tidak terjadi duplikasi saat server restart
		fmt.Println("[Seeder] ℹ️ Akun Super Admin sudah tersedia, melewati proses seeding.")
	}

}
