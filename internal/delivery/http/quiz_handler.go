package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type QuizHandler struct {
	quizUsecase     domain.QuizUsecase
	questionUsecase domain.QuestionUsecase
}

func NewQuizHandler(r *gin.Engine, qu domain.QuizUsecase, qnu domain.QuestionUsecase) {
	handler := &QuizHandler{
		quizUsecase:     qu,
		questionUsecase: qnu,
	}

	api := r.Group("/api")
	{
		// Quiz Management
		api.POST("/modules/:module_id/quizzes", handler.CreateQuiz)
		api.GET("/modules/:module_id/quizzes", handler.GetQuizzesByModule)
		api.GET("/quizzes/:quiz_id", handler.GetQuizByID)
		api.PUT("/quizzes/:quiz_id", handler.UpdateQuiz)
		api.DELETE("/quizzes/:quiz_id", handler.DeleteQuiz)

		// Question Management (Nested under Quiz)
		api.POST("/quizzes/:quiz_id/questions", handler.CreateQuestion)
		api.GET("/quizzes/:quiz_id/questions", handler.GetQuestionsByQuiz)
		api.GET("/questions/:question_id", handler.GetQuestionByID)
		api.PUT("/questions/:question_id", handler.UpdateQuestion)
		api.DELETE("/questions/:question_id", handler.DeleteQuestion)
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

// ... (UpdateQuiz & DeleteQuiz mengikuti pola yang sama) ...

// --- Question Handlers ---

type QuestionRequest struct {
	Type          domain.QuestionType `json:"type" binding:"required"`
	Text          string              `json:"text" binding:"required"`
	Options       datatypes.JSON      `json:"options"`
	CorrectAnswer datatypes.JSON      `json:"correct_answer" binding:"required"`
	Points        int                 `json:"points"`
	Explanation   string              `json:"explanation"`
}

func (h *QuizHandler) CreateQuestion(c *gin.Context) {
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

func (h *QuizHandler) GetQuestionsByQuiz(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("quiz_id"), 10, 32)
	questions, err := h.questionUsecase.GetQuestionsByQuizID(uint(quizID), true) // true untuk randomize
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

func (h *QuizHandler) GetQuestionByID(c *gin.Context) {
	qID, _ := strconv.ParseUint(c.Param("question_id"), 10, 32)
	question, err := h.questionUsecase.GetQuestionByID(uint(qID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (h *QuizHandler) UpdateQuestion(c *gin.Context) {
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

func (h *QuizHandler) DeleteQuestion(c *gin.Context) {
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
