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
}

func NewCourseHandler(r *gin.Engine, cu domain.CourseUsecase, er domain.EnrollmentRepository) {
	handler := &CourseHandler{
		courseUsecase: cu,
	}

	// Akses untuk semua yang login (Siswa, Tutor, Admin)
	courseGeneral := r.Group("/api/courses")
	courseGeneral.Use(middleware.RequireAuth())
	{
		courseGeneral.GET("", handler.GetAll) // List Courses
	}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	courseProtected := courseGeneral.Group("/:course_id")
	courseProtected.Use(middleware.RequireCourseAccess(er))
	{
		courseProtected.GET("", handler.GetByID) // Get Detail by ID
	}

	// Akses Hanya untuk Tutor dan Admin
	coursePrivate := courseGeneral.Group("")
	coursePrivate.Use(middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		coursePrivate.POST("", handler.Create)              // Create
		coursePrivate.PUT("/:course_id", handler.Update)    // Update
		coursePrivate.DELETE("/:course_id", handler.Delete) // Delete
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
