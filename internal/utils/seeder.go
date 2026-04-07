package utils

import (
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "admin@gmail.com").Count(&count)

	if count == 0 {
		fmt.Println("[Seeder] Mendeteksi database baru. Menyiapkan akun Super Admin...")

		hashedPassword, err := utils.HashPassword("admin123")
		if err != nil {
			fmt.Println("❌ [Seeder] Gagal hash password Super Admin:", err)
		}

		admin := domain.User{
			Name:     "Admin Ganteng",
			Email:    "admin@gmail.com",
			Password: string(hashedPassword),
			Role:     "Admin",
		}

		if err := db.Create(&admin).Error; err != nil {
			fmt.Println("❌ [Seeder] Gagal membuat akun Super Admin:", err)
		} else {
			fmt.Println("[Seeder] ✅ Berhasil membuat Super Admin! (Email: admin@gmail.com | Pass: admin123)")
		}
	} else {
		fmt.Println("[Seeder] ℹ️ Akun Super Admin sudah tersedia, melewati proses seeding.")
	}

}
