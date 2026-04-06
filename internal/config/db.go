package config

import (
	"fmt"
	"log"
	"os"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// Migrasi otomatis
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Course{},
		&domain.Module{},
		&domain.Lesson{},
		&domain.Quiz{},
		&domain.Question{},
		&domain.QuizAttempt{},
		&domain.Exam{},
		&domain.ExamQuestion{},
		&domain.ExamAttempt{},
		&domain.Enrollment{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi:", err)
	}
	return db
}
