package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type EnrollmentHandler struct {
	enrollUsecase domain.EnrollmentUsecase
}

func NewEnrollmentHandler(r *gin.Engine, eu domain.EnrollmentUsecase, er domain.EnrollmentRepository) {
	handler := &EnrollmentHandler{
		enrollUsecase: eu,
	}

	// Akses Hanya untuk Tutor dan Admin
	coursePrivate := r.Group("/api/courses")
	coursePrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		coursePrivate.GET("/:course_id/enrollments", handler.GetEnrollments)
		coursePrivate.POST("/:course_id/enrollments/:user_id", handler.EnrollStudent)
		coursePrivate.DELETE("/:course_id/enrollments/:user_id", handler.UnenrollStudent)
	}
}

func (h *EnrollmentHandler) GetEnrollments(c *gin.Context) {
	idParam := c.Param("course_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	enrollments, err := h.enrollUsecase.GetEnrolledStudents(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": enrollments})
}

func (h *EnrollmentHandler) EnrollStudent(c *gin.Context) {
	courseIdParam := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	studentIdParam := c.Param("user_id")
	studentID, err := strconv.ParseUint(studentIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.enrollUsecase.EnrollStudent(uint(courseID), uint(studentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil didaftarkan ke mata pelajaran"})
}

func (h *EnrollmentHandler) UnenrollStudent(c *gin.Context) {
	courseIdParam := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	studentIdParam := c.Param("user_id")
	studentID, err := strconv.ParseUint(studentIdParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.enrollUsecase.UnenrollStudent(uint(courseID), uint(studentID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil dihapus dari mata pelajaran"})
}
