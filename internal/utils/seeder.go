package utils

import (
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("username = ?", "admin").Count(&count)

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if count == 0 {
		adminUser := domain.User{
			Email:    "admin@gmail.com",
			Password: string(hash),
			Role:     "admin123",
		}
		if err := db.Create(&adminUser).Error; err != nil {
			fmt.Println("❌ Gagal membuat akun admin:", err)
		} else {
			fmt.Println("✅ SEEDER: Akun Admin berhasil dibuat (admin@gmail.com / admin123)!")
		}
	} else {
		fmt.Println("✅ SEEDER: Akun Admin sudah eksis, melewati proses seeding.")
	}
}
