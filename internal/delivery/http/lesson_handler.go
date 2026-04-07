package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	lessonUsecase domain.LessonUsecase
}

func NewLessonHandler(r *gin.Engine, mu domain.LessonUsecase, er domain.EnrollmentRepository) {
	handler := &LessonHandler{lessonUsecase: mu}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	lessonProtected := r.Group("/api")
	lessonProtected.Use(middleware.RequireAuth(), middleware.RequireModuleAccess(er))
	{
		lessonProtected.GET("/modules/:module_id/lessons", handler.GetByModuleID)
	}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	lessonProtected2 := r.Group("/api")
	lessonProtected2.Use(middleware.RequireAuth(), middleware.RequireLessonAccess(er))
	{
		lessonProtected2.GET("/lessons/:lesson_id", handler.GetByID)
	}
	// Akses Hanya untuk Tutor dan Admin
	lessonPrivate := r.Group("/api")
	lessonPrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		lessonPrivate.POST("/modules/:module_id/lessons", handler.Create)
		lessonPrivate.PUT("/lessons/:lesson_id", handler.Update)
		lessonPrivate.DELETE("/lessons/:lesson_id", handler.Delete)
	}
}

type LessonRequest struct {
	Title      string `json:"title" binding:"required"`
	LessonType string `json:"lesson_type" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Sequence   int    `json:"sequence"`
}

func (h *LessonHandler) Create(c *gin.Context) {
	courseID, _ := strconv.ParseUint(c.Param("module_id"), 10, 32)
	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module, err := h.lessonUsecase.CreateLesson(uint(courseID), req.Title, domain.LessonType(req.LessonType), req.Content, req.Sequence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": module})
}

func (h *LessonHandler) GetByModuleID(c *gin.Context) {
	moduleId := c.Param("module_id")
	id, err := strconv.ParseUint(moduleId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	courses, err := h.lessonUsecase.GetLessonsByModule(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *LessonHandler) GetByID(c *gin.Context) {
	idParam := c.Param("lesson_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	course, err := h.lessonUsecase.GetLessonByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (h *LessonHandler) Update(c *gin.Context) {
	idParam := c.Param("lesson_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req LessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	course, err := h.lessonUsecase.UpdateLesson(uint(id), req.Title, domain.LessonType(req.LessonType), req.Content, req.Sequence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mata pelajaran berhasil diperbarui",
		"data":    course,
	})
}

func (h *LessonHandler) Delete(c *gin.Context) {
	idParam := c.Param("lesson_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.lessonUsecase.DeleteLesson(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}
