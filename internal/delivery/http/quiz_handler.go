package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type QuizHandler struct {
	quizUsecase domain.QuizUsecase
}

func NewQuizHandler(r *gin.Engine, qu domain.QuizUsecase, er domain.EnrollmentRepository) {
	handler := &QuizHandler{
		quizUsecase: qu,
	}

	quizProtected := r.Group("/api")
	quizProtected.Use(middleware.RequireAuth(), middleware.RequireModuleAccess(er))
	{
		quizProtected.GET("/modules/:module_id/quizzes", handler.GetQuizzesByModule)
	}
	quizProtected2 := r.Group("/api")
	quizProtected2.Use(middleware.RequireAuth(), middleware.RequireQuizAccess(er))
	{
		quizProtected2.GET("/quizzes/:quiz_id", handler.GetQuizByID)
	}

	quizPrivate := r.Group("/api")
	quizPrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		// Quiz Management
		quizPrivate.POST("/modules/:module_id/quizzes", handler.CreateQuiz)
		quizPrivate.PUT("/quizzes/:quiz_id", handler.UpdateQuiz)
		quizPrivate.DELETE("/quizzes/:quiz_id", handler.DeleteQuiz)
		// Question Management (Nested under Quiz)
		quizPrivate.POST("/quizzes/:quiz_id/questions/generate", handler.GenerateQuizQuestionsWithAI)
	}
}

// --- Quiz Handlers ---

type QuizRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	TimeLimit    int    `json:"time_limit"`
	PassingScore int    `json:"passing_score"`
	IsRandomized bool   `json:"is_randomized"`
	MaxAttempts  int    `json:"max_attempts"`
}

func (h *QuizHandler) CreateQuiz(c *gin.Context) {
	moduleID, _ := strconv.ParseUint(c.Param("module_id"), 10, 32)
	var req QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	quiz, err := h.quizUsecase.CreateQuiz(uint(moduleID), req.Title, req.Description, req.TimeLimit, req.PassingScore, req.IsRandomized, req.MaxAttempts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": quiz})
}

func (h *QuizHandler) GetQuizzesByModule(c *gin.Context) {
	moduleID, _ := strconv.ParseUint(c.Param("module_id"), 10, 32)
	quizzes, err := h.quizUsecase.GetQuizzesByModule(uint(moduleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quizzes})
}

func (h *QuizHandler) GetQuizByID(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)
	quiz, err := h.quizUsecase.GetQuizByID(uint(quizID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kuis tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": quiz})
}

func (h *QuizHandler) UpdateQuiz(c *gin.Context) {
	idParam := c.Param("quiz_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	quiz, err := h.quizUsecase.UpdateQuiz(uint(id), req.Title, req.Description, req.TimeLimit, req.PassingScore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kuis berhasil diperbarui",
		"data":    quiz,
	})
}

func (h *QuizHandler) DeleteQuiz(c *gin.Context) {
	idParam := c.Param("quiz_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.quizUsecase.DeleteQuiz(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kuis berhasil dihapus"})

}

func (h *QuizHandler) GenerateQuizQuestionsWithAI(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)

	var req AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	questions, err := h.quizUsecase.GenerateQuizQuestionsWithAI(uint(quizID), req.Topic, req.QType, req.Count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Berhasil menggenerate soal menggunakan AI",
		"data":    questions,
	})
}

func (h *QuizHandler) StartAttempt(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)
	userIDVal, _ := c.Get("user_id") // Pastikan route ini diproteksi AuthMiddleware
	userID := uint(userIDVal.(float64))

	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token Ujian (CBT Token) diperlukan"})
		return
	}

	attempt, questions, err := h.quizUsecase.StartAttempt(uint(quizID), userID)
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
