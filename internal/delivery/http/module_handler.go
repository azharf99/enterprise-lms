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

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	moduleProtected := r.Group("/api")
	moduleProtected.Use(middleware.RequireAuth(), middleware.RequireModuleAccess(er))
	{
		moduleProtected.GET("/modules/:module_id", handler.GetByID)
	}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	moduleProtected2 := r.Group("/api")
	moduleProtected2.Use(middleware.RequireAuth(), middleware.RequireCourseAccess(er))
	{
		moduleProtected2.GET("/courses/:course_id/modules", handler.GetByCourse)
	}

	// Akses untuk Siswa yang terdaftar di mata pelajaran, Tutor dan Admin
	// Akses Hanya untuk Tutor dan Admin
	modulePrivate := r.Group("/api")
	modulePrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		modulePrivate.POST("/courses/:course_id/modules", handler.Create)
		modulePrivate.PUT("/modules/:module_id", handler.Update)
		modulePrivate.DELETE("/modules/:module_id", handler.Delete)

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
