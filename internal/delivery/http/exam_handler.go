package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	examUsecase domain.ExamUsecase
}

func NewExamHandler(r *gin.Engine, eu domain.ExamUsecase, er domain.EnrollmentRepository) {
	handler := &ExamHandler{
		examUsecase: eu,
	}
	// Endpoint Manajemen Ujian (Admin/Tutor)
	mgmt := r.Group("/api/courses/:course_id/exams")
	mgmt.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		mgmt.POST("", handler.CreateExam)
		mgmt.GET("", handler.GetExamsByCourseID)
	}

	examMgmt := r.Group("/api/exams")
	examMgmt.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		examMgmt.GET("/:exam_id", handler.GetExamByID)
		examMgmt.PUT("/:exam_id", handler.UpdateExam)
		examMgmt.DELETE("/:exam_id", handler.DeleteExam)
		examMgmt.PATCH("/:exam_id/token", handler.GenerateToken)
		examMgmt.POST("/:exam_id/questions/generate", handler.GenerateQuestionsAI)
	}
}

// --- Implementasi CreateExam ---
func (h *ExamHandler) CreateExam(c *gin.Context) {
	// Ambil course_id dari parameter URL (/api/courses/:course_id/exams)
	courseIDParam := c.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Course tidak valid"})
		return
	}

	var req domain.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	// Panggil Usecase
	exam, err := h.examUsecase.CreateExam(uint(courseID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ujian berhasil dibuat",
		"data":    exam,
	})
}

func (h *ExamHandler) GetExamsByCourseID(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := strconv.ParseUint(courseId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	exams, err := h.examUsecase.GetExamsByCourseID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": exams})
}

func (h *ExamHandler) GetExamByID(c *gin.Context) {
	examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	exam, err := h.examUsecase.GetExamByID(uint(examID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": exam})
}

func (h *ExamHandler) UpdateExam(c *gin.Context) {
	idParam := c.Param("exam_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req domain.CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	exam, err := h.examUsecase.UpdateExam(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kuis berhasil diperbarui",
		"data":    exam,
	})
}


func (h *ExamHandler) DeleteExam(c *gin.Context) {
	idParam := c.Param("exam_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.examUsecase.DeleteExam(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ujian berhasil dihapus"})
}

// --- Implementasi GenerateToken ---
func (h *ExamHandler) GenerateToken(c *gin.Context) {
	// Ambil exam_id dari parameter URL (/api/exams/:exam_id/token)
	examIDParam := c.Param("exam_id")
	examID, err := strconv.ParseUint(examIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Ujian tidak valid"})
		return
	}

	// Panggil Usecase untuk mengacak token baru
	token, err := h.examUsecase.GenerateCBTToken(uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token CBT baru berhasil dibuat",
		"token":   token,
	})
}

func (h *ExamHandler) GenerateQuestionsAI(c *gin.Context) {
	examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)

	var req domain.AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := h.examUsecase.GenerateExamQuestionsWithAI(uint(examID), req.Topic, req.QType, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menggenerate soal menggunakan AI",
		"data":    questions,
	})
}

func (h *ExamHandler) StartAttempt(c *gin.Context) {
	examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	userIDVal, _ := c.Get("user_id") // Pastikan route ini diproteksi AuthMiddleware
	userID := uint(userIDVal.(float64))

	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token Ujian (CBT Token) diperlukan"})
		return
	}

	attempt, questions, err := h.examUsecase.StartExamAttempt(uint(examID), userID, req.Token)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Ujian dimulai",
		"attempt":   attempt,
		"questions": questions,
	})
}
