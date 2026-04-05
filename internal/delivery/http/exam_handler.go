package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type ExamHandler struct {
	examUsecase domain.ExamUsecase
}

func NewExamHandler(r *gin.Engine, eu domain.ExamUsecase) {
	handler := &ExamHandler{examUsecase: eu}

	// Endpoint Manajemen Ujian (Admin/Tutor)
	mgmt := r.Group("/api/courses/:course_id/exams")
	{
		mgmt.POST("", handler.CreateExam)
	}

	api := r.Group("/api/exams/:exam_id")
	{
		api.PATCH("/token", handler.GenerateToken)
		api.POST("/generate-ai", handler.GenerateQuestionsAI)
		api.POST("/attempts", handler.StartAttempt)
	}
}

// Format Request untuk Generate AI
type AIGenerateRequest struct {
	Topic string `json:"topic" binding:"required"`
	QType string `json:"q_type" binding:"required"`
	Count int    `json:"count" binding:"required,min=1,max=25"`
}

// --- Struct untuk Request Body ---
type CreateExamRequest struct {
	Title        string     `json:"title" binding:"required"`
	ExamType     string     `json:"exam_type" binding:"required"` // Misalnya: "PTS", "PAS"
	Description  string     `json:"description"`
	TimeLimit    int        `json:"time_limit"`
	PassingScore int        `json:"passing_score"`
	StartTime    *time.Time `json:"start_time"` // Format JSON harus RFC3339, misal: "2026-10-01T08:00:00Z"
	EndTime      *time.Time `json:"end_time"`
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

	var req CreateExamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid: " + err.Error()})
		return
	}

	// Panggil Usecase
	exam, err := h.examUsecase.CreateExam(
		uint(courseID), req.Title, req.ExamType, req.Description,
		req.TimeLimit, req.PassingScore, req.StartTime, req.EndTime,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ujian berhasil dibuat",
		"data":    exam,
	})
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

	var req AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := h.examUsecase.GenerateQuestionsWithAI(uint(examID), req.Topic, req.QType, req.Count)
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

	attempt, questions, err := h.examUsecase.StartAttempt(uint(examID), userID, req.Token)
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

// ... (Tambahkan fungsi pendukung lainnya sesuai kebutuhan) ...
