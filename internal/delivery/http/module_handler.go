package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type ModuleHandler struct {
	moduleUsecase domain.ModuleUsecase
}

func NewModuleHandler(r *gin.Engine, mu domain.ModuleUsecase, er domain.EnrollmentRepository) {
	handler := &ModuleHandler{moduleUsecase: mu}

	lessonAuth := r.Group("/api")
	lessonAuth.Use(middleware.RequireAuth(), middleware.RequireCourseAccess(er))
	// Contoh: /api/courses/1/modules
	lessonAuth.GET("/courses/:course_id/modules", handler.GetByCourse)

	lessonMgmt := r.Group("/api")
	lessonMgmt.Use(middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		// Contoh: /api/courses/1/modules
		lessonMgmt.POST("/courses/:course_id/modules", handler.Create)
		lessonMgmt.GET("/modules/:module_id", handler.GetByID)
		lessonMgmt.PUT("/modules/:module_id", handler.Update)
		lessonMgmt.DELETE("/modules/:module_id", handler.Delete)
	}
}

type ModuleRequest struct {
	Title    string `json:"title" binding:"required"`
	Sequence int    `json:"sequence"`
}

func (h *ModuleHandler) Create(c *gin.Context) {
	courseID, _ := strconv.ParseUint(c.Param("course_id"), 10, 32)
	var req ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module, err := h.moduleUsecase.CreateModule(uint(courseID), req.Title, req.Sequence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": module})
}

func (h *ModuleHandler) GetByCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := strconv.ParseUint(courseId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	courses, err := h.moduleUsecase.GetModulesByCourse(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": courses})
}

func (h *ModuleHandler) GetByID(c *gin.Context) {
	idParam := c.Param("module_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	course, err := h.moduleUsecase.GetModuleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": course})
}

func (h *ModuleHandler) Update(c *gin.Context) {
	idParam := c.Param("module_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req ModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	course, err := h.moduleUsecase.UpdateModule(uint(id), req.Title, req.Sequence)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mata pelajaran berhasil diperbarui",
		"data":    course,
	})
}

func (h *ModuleHandler) Delete(c *gin.Context) {
	idParam := c.Param("module_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.moduleUsecase.DeleteModule(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}
