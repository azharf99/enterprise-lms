package middleware

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

// 1. RequireCourseAccess: Untuk rute /api/courses/:course_id/...
func RequireCourseAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" { c.Next(); return }

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))

		courseID, _ := strconv.ParseUint(c.Param("course_id"), 10, 32)

		isEnrolled, _ := enrollmentRepo.CheckEnrollment(uint(courseID), userID)
		if !isEnrolled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak terdaftar dalam kursus ini."})
			return
		}
		c.Next()
	}
}

// 2. RequireModuleAccess: Untuk rute /api/modules/:module_id/...
func RequireModuleAccess(moduleRepo domain.ModuleRepository, enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" { c.Next(); return }

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		moduleID, _ := strconv.ParseUint(c.Param("module_id"), 10, 32)

		module, err := moduleRepo.GetByID(uint(moduleID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Modul tidak ditemukan"})
			return
		}

		isEnrolled, _ := enrollmentRepo.CheckEnrollment(module.CourseID, userID)
		if !isEnrolled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak terdaftar dalam kursus modul ini."})
			return
		}
		c.Next()
	}
}

// 3. RequireQuizAccess: Untuk rute /api/quizzes/:quiz_id/...
func RequireQuizAccess(quizRepo domain.QuizRepository, moduleRepo domain.ModuleRepository, enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" { c.Next(); return }

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)

		quiz, err := quizRepo.GetByID(uint(quizID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
			return
		}

		module, err := moduleRepo.GetByID(quiz.ModuleID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Modul kuis tidak ditemukan"})
			return
		}

		isEnrolled, _ := enrollmentRepo.CheckEnrollment(module.CourseID, userID)
		if !isEnrolled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses kuis ini."})
			return
		}
		c.Next()
	}
}

// 4. RequireExamAccess: Untuk rute /api/exams/:exam_id/...
func RequireExamAccess(examRepo domain.ExamRepository, enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" { c.Next(); return }

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)

		exam, err := examRepo.GetByID(uint(examID))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
			return
		}

		// Karena di Fase 4 kita mendesain Exam terhubung langsung ke Course (bukan Module)
		isEnrolled, _ := enrollmentRepo.CheckEnrollment(exam.CourseID, userID)
		if !isEnrolled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses ujian ini."})
			return
		}
		c.Next()
	}
}