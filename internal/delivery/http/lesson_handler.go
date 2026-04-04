package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	lessonUsecase domain.LessonUsecase
}

func NewLessonHandler(r *gin.Engine, mu domain.LessonUsecase) {
	handler := &LessonHandler{lessonUsecase: mu}

	// API bersarang (nested) lebih disarankan untuk REST yang baik
	// Contoh: /api/courses/1/lessons
	api := r.Group("/api")
	{
		api.POST("/modules/:moduleId/lessons", handler.Create)
		api.GET("/modules/:moduleId/lessons", handler.GetByModuleID)
		api.GET("/lessons/:id", handler.GetByID)
		api.PUT("/lessons/:id", handler.Update)
		api.DELETE("/lessons/:id", handler.Delete)
	}
}

type LessonRequest struct {
	Title      string `json:"title" binding:"required"`
	LessonType string `json:"lessonType" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Sequence   int    `json:"sequence"`
}

func (h *LessonHandler) Create(c *gin.Context) {
	courseID, _ := strconv.ParseUint(c.Param("moduleId"), 10, 32)
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
	moduleId := c.Param("moduleId")
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
	idParam := c.Param("id")
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
	idParam := c.Param("id")
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
	idParam := c.Param("id")
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
