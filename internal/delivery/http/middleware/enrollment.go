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
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

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

// 2. RequireModuleAccess: Versi Teroptimasi (Single Query)
func RequireModuleAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		moduleID, _ := strconv.ParseUint(c.Param("module_id"), 10, 32)

		hasAccess, err := enrollmentRepo.CheckModuleAccess(uint(moduleID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses modul ini."})
			return
		}
		c.Next()
	}
}

// 4. RequireExamAccess: Versi Teroptimasi (Single Query)
func RequireExamAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)

		hasAccess, err := enrollmentRepo.CheckExamAccess(uint(examID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses ujian ini."})
			return
		}
		c.Next()
	}
}

// 4. RequireExamAttemptAccess: Versi Teroptimasi (Single Query)
func RequireExamAttemptAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		attemptID, _ := strconv.ParseUint(c.Param("attempt_id"), 10, 32)

		hasAccess, err := enrollmentRepo.CheckExamAttemptAccess(uint(attemptID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses ujian ini."})
			return
		}
		c.Next()
	}
}

// 4. RequireQuizAttemptAccess: Versi Teroptimasi (Single Query)
func RequireQuizAttemptAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		attemptID, _ := strconv.ParseUint(c.Param("attempt_id"), 10, 32)

		hasAccess, err := enrollmentRepo.CheckQuizAttemptAccess(uint(attemptID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak. Anda tidak berhak mengakses ujian ini."})
			return
		}
		c.Next()
	}
}

func RequireQuizAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)

		// HANYA 1 QUERY SAJA!
		hasAccess, err := enrollmentRepo.CheckQuizAccess(uint(quizID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak."})
			return
		}

		c.Next()
	}
}

func RequireQuestionAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		questionID, _ := strconv.ParseUint(c.Param("question_id"), 10, 32)

		// HANYA 1 QUERY SAJA!
		hasAccess, err := enrollmentRepo.CheckQuestionAccess(uint(questionID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak."})
			return
		}

		c.Next()
	}
}

func RequireLessonAccess(enrollmentRepo domain.EnrollmentRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		if userRole == "Admin" || userRole == "Tutor" {
			c.Next()
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID := uint(userIDVal.(float64))
		lessonID, _ := strconv.ParseUint(c.Param("lesson_id"), 10, 32)

		// HANYA 1 QUERY SAJA!
		hasAccess, err := enrollmentRepo.CheckLessonAccess(uint(lessonID), userID)
		if err != nil || !hasAccess {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Akses Ditolak."})
			return
		}

		c.Next()
	}
}
