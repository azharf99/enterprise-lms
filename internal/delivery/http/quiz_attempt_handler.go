package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type AttemptHandler struct {
	quizUsecase domain.QuizUsecase
}

func NewAttemptHandler(r *gin.Engine, qu domain.QuizUsecase, er domain.EnrollmentRepository) {
	handler := &AttemptHandler{quizUsecase: qu}

	examProtected1 := r.Group("/api")
	examProtected1.Use(middleware.RequireAuth(), middleware.RequireQuizAccess(er))
	{
		// Memulai kuis
		examProtected1.POST("/quizzes/:quiz_id/attempts", handler.StartAttempt)
	}
	examProtected2 := r.Group("/api")
	examProtected2.Use(middleware.RequireAuth(), middleware.RequireQuizAttemptAccess(er))
	{
		// Mengirimkan jawaban kuis
		examProtected2.POST("/attempts/:attempt_id/submit", handler.SubmitAttempt)
	}
}

func (h *AttemptHandler) StartAttempt(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)

	// c.Get("user_id") mengembalikan interface{}, kita konversi ke uint
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDVal.(float64)) // JWT parsing biasanya menghasilkan float64

	response, err := h.quizUsecase.StartAttempt(uint(quizID), userID, "in_progress") // Kita set status awal "in_progress"
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Kuis dimulai",
		"attempt":   response.Attempt,
		"questions": response.Questions,
	})
}

func (h *AttemptHandler) SubmitAttempt(c *gin.Context) {
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

	result, err := h.quizUsecase.SubmitAttempt(uint(attemptID), answersJSON)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kuis berhasil dikirim",
		"score":   result.Score,
		"passed":  result.Passed,
		"status":  "Completed",
	})
}
