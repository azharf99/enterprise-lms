package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type QuizQuestionHandler struct {
	questionUsecase domain.QuizQuestionUsecase
}

func NewQuizQuestionHandler(r *gin.Engine, qnu domain.QuizQuestionUsecase, er domain.EnrollmentRepository) {
	handler := &QuizQuestionHandler{
		questionUsecase: qnu,
	}

	quizProtected := r.Group("/api")
	quizProtected.Use(middleware.RequireAuth(), middleware.RequireQuizAccess(er))
	{
		quizProtected.GET("/quizzes/:quiz_id/questions", handler.GetQuestionsByQuiz)
	}

	quizProtected3 := r.Group("/api")
	quizProtected3.Use(middleware.RequireAuth(), middleware.RequireQuestionAccess(er))
	{
		quizProtected3.GET("/questions/:question_id", handler.GetQuestionByID)
	}

	quizPrivate := r.Group("/api")
	quizPrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		// Question Management (Nested under Quiz)
		quizPrivate.POST("/quizzes/:quiz_id/questions", handler.CreateQuestion)
		quizPrivate.PUT("/questions/:question_id", handler.UpdateQuestion)
		quizPrivate.DELETE("/questions/:question_id", handler.DeleteQuestion)
	}
}

type QuestionRequest struct {
	Type          domain.QuestionType `json:"type" binding:"required"`
	Text          string              `json:"text" binding:"required"`
	Options       datatypes.JSON      `json:"options"`
	CorrectAnswer datatypes.JSON      `json:"correct_answer" binding:"required"`
	Points        int                 `json:"points"`
	Explanation   string              `json:"explanation"`
}

func (h *QuizQuestionHandler) CreateQuestion(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)
	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.questionUsecase.CreateQuestion(uint(quizID), req.Type, req.Text, req.Options, req.CorrectAnswer, req.Points, req.Explanation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": question})
}

func (h *QuizQuestionHandler) GetQuestionsByQuiz(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)
	questions, err := h.questionUsecase.GetQuestionsByQuizID(uint(quizID), true) // true untuk randomize
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

func (h *QuizQuestionHandler) GetQuestionByID(c *gin.Context) {
	qID, _ := strconv.ParseUint(c.Param("question_id"), 10, 32)
	question, err := h.questionUsecase.GetQuestionByID(uint(qID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (h *QuizQuestionHandler) UpdateQuestion(c *gin.Context) {
	idParam := c.Param("question_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	question, err := h.questionUsecase.UpdateQuestion(uint(id), req.Type, req.Text, req.Options, req.CorrectAnswer, req.Points, req.Explanation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Soal berhasil diperbarui",
		"data":    question,
	})
}

func (h *QuizQuestionHandler) DeleteQuestion(c *gin.Context) {
	idParam := c.Param("question_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.questionUsecase.DeleteQuestion(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dihapus"})
}
