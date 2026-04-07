package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type ExamAttemptHandler struct {
	examUsecase domain.ExamUsecase
}

func NewExamAttemptHandler(r *gin.Engine, eu domain.ExamUsecase, er domain.EnrollmentRepository) {
	handler := &ExamAttemptHandler{
		examUsecase: eu,
	}

	examProtected1 := r.Group("/api/exams/:exam_id")
	examProtected1.Use(middleware.RequireAuth(), middleware.RequireExamAccess(er))
	{
		examProtected1.POST("/attempts", handler.StartExamAttempt)
	}
	examProtected2 := r.Group("/api")
	examProtected2.Use(middleware.RequireAuth(), middleware.RequireExamAttemptAccess(er))
	{
		// Mengirimkan jawaban kuis
		examProtected2.POST("/exam-attempts/:attempt_id/submit", handler.SubmitExamAttempt)
	}
}

func (h *ExamAttemptHandler) StartExamAttempt(c *gin.Context) {
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

func (h *ExamAttemptHandler) SubmitExamAttempt(c *gin.Context) {
	attemptID, _ := strconv.ParseUint(c.Param("attempt_id"), 10, 32)

	var req struct {
		Answers interface{} `json:"answers" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jawaban diperlukan"})
		return
	}

	// Konversi input answers ke JSONB format
	answersJSON, _ := json.Marshal(req.Answers)

	result, err := h.examUsecase.SubmitExamAttempt(uint(attemptID), answersJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kuis berhasil dikirim",
		"score":   result.Score,
		"status":  "Completed",
	})
}
