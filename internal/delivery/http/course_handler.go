package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	courseUsecase domain.CourseUsecase
	enrollUsecase domain.EnrollmentUsecase
}

func NewCourseHandler(r *gin.Engine, cu domain.CourseUsecase, eu domain.EnrollmentUsecase, er domain.EnrollmentRepository) {
	handler := &CourseHandler{
		courseUsecase: cu,
		enrollUsecase: eu,
	}

	// Akses untuk semua yang login (Siswa, Tutor, Admin)
	courseGeneral := r.Group("/api/courses")
	courseGeneral.Use(middleware.RequireAuth())
	{
		courseGeneral.GET("", handler.GetAll)
	}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	courseProtected := courseGeneral.Group("/:course_id")
	courseProtected.Use(middleware.RequireCourseAccess(er))
	{
		courseProtected.GET("", handler.GetByID)
	}

	// Akses Hanya untuk Tutor dan Admin
	coursePrivate := courseGeneral.Group("")
	coursePrivate.Use(middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		coursePrivate.POST("", handler.Create)
		coursePrivate.PUT("/:course_id", handler.Update)
		coursePrivate.GET("/:course_id/enrollments", handler.GetEnrollments)
		coursePrivate.POST("/:course_id/enrollments/:user_id", handler.EnrollStudent)
		coursePrivate.DELETE("/:course_id/enrollments/:user_id", handler.UnenrollStudent)
		coursePrivate.DELETE("/:course_id", handler.Delete)
	}
}

type CourseRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	TutorIDs    []uint `json:"tutor_ids"`
}

// ... (handler Create, GetAll, GetByID tetap sama) ...
func (h *CourseHandler) Create(c *gin.Context) {
	var req CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}
	course, err := h.courseUsecase.CreateCourse(req.Title, req.Description, req.TutorIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Mata pelajaran berhasil dibuat", "data": course})
}

func (h *CourseHandler) GetAll(c *gin.Context) {
	courses, err := h.courseUsecase.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	idParam := c.Param("course_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	course, err := h.courseUsecase.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (h *CourseHandler) Update(c *gin.Context) {
	idParam := c.Param("course_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req CourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	course, err := h.courseUsecase.UpdateCourse(uint(id), req.Title, req.Description, req.TutorIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mata pelajaran berhasil diperbarui",
		"data":    course,
	})
}

func (h *CourseHandler) Delete(c *gin.Context) {
	idParam := c.Param("course_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.courseUsecase.DeleteCourse(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}

func (h *CourseHandler) GetEnrollments(c *gin.Context) {
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

func (h *CourseHandler) EnrollStudent(c *gin.Context) {
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

func (h *CourseHandler) UnenrollStudent(c *gin.Context) {
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
