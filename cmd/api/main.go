package main

import (
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/config"
	"github.com/azharf99/enterprise-lms/internal/delivery/http"
	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/repository/postgres"
	"github.com/azharf99/enterprise-lms/internal/usecase"
	"github.com/azharf99/enterprise-lms/internal/utils"
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
	utils.SeedAdmin(db)

	// 2. Inisialisasi Repository
	userRepo := postgres.NewUserRepository(db)
	courseRepo := postgres.NewCourseRepository(db)
	moduleRepo := postgres.NewModuleRepository(db)
	lessonRepo := postgres.NewLessonRepository(db)
	quizRepo := postgres.NewQuizRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	attemptRepo := postgres.NewQuizAttemptRepository(db)
	examRepo := postgres.NewExamRepository(db)
	examQuestionRepo := postgres.NewExamQuestionRepository(db)
	examAttemptRepo := postgres.NewExamAttemptRepository(db)
	enrollmentRepo := postgres.NewEnrollmentRepository(db)

	// 3. Inisialisasi Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	courseUsecase := usecase.NewCourseUsecase(courseRepo)
	moduleUsecase := usecase.NewModuleUsecase(moduleRepo)
	lessonUsecase := usecase.NewLessonUsecase(lessonRepo)
	questionUsecase := usecase.NewQuestionUsecase(questionRepo)
	quizUsecase := usecase.NewQuizUsecase(quizRepo, attemptRepo)
	examUsecase := usecase.NewExamUsecase(examRepo, examQuestionRepo, examAttemptRepo)
	analyticsUsecase := usecase.NewAnalyticsUsecase(examRepo, examAttemptRepo)
	enrollmentUsecase := usecase.NewEnrollmentUsecase(enrollmentRepo)
	// 4. Inisialisasi Router & Handler
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	// Pasang middleware keamanan global
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.SetupCORS())
	r.Use(middleware.RateLimiter())
	http.NewUserHandler(r, userUsecase)
	protedcted := r.Group("/")
	protedcted.Use(middleware.RequireAuth())
	{
		http.NewCourseHandler(r, courseUsecase, enrollmentUsecase, enrollmentRepo)
		http.NewModuleHandler(r, moduleUsecase, enrollmentRepo)
		http.NewLessonHandler(r, lessonUsecase, enrollmentRepo)
		http.NewQuizHandler(r, quizUsecase, questionUsecase, enrollmentRepo)
		http.NewAttemptHandler(r, quizUsecase)
		http.NewExamHandler(r, examUsecase)
		http.NewAnalyticsHandler(r, analyticsUsecase)
	}

	// 5. Jalankan Server
	r.Run(":8080")
}
