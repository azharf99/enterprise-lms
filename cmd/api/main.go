package main

import (
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/config"
	"github.com/azharf99/enterprise-lms/internal/delivery/http"
	"github.com/azharf99/enterprise-lms/internal/repository/postgres"
	"github.com/azharf99/enterprise-lms/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi ENV dan Database
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file. Mengalihkan dengan menggunakan ENV dari docker.")
	}
	db := config.ConnectDatabase()

	// 2. Inisialisasi Repository
	userRepo := postgres.NewUserRepository(db)

	// 3. Inisialisasi Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 4. Inisialisasi Router & Handler
	r := gin.Default()
	http.NewUserHandler(r, userUsecase)

	// 5. Jalankan Server
	r.Run(":8081")
}
