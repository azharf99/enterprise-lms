package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ExamQuestionHandler struct {
	examQuestionUsecase domain.ExamQuestionUsecase
}

func NewExamQuestionHandler(r *gin.Engine, qnu domain.ExamQuestionUsecase, er domain.EnrollmentRepository) {
	handler := &ExamQuestionHandler{
		examQuestionUsecase: qnu,
	}

	quizProtected := r.Group("/api")
	quizProtected.Use(middleware.RequireAuth(), middleware.RequireQuizAccess(er))
	{
		quizProtected.GET("/exams/:exam_id/questions", handler.GetExamQuestionsByExamID)
	}

	quizProtected3 := r.Group("/api")
	quizProtected3.Use(middleware.RequireAuth(), middleware.RequireQuestionAccess(er))
	{
		quizProtected3.GET("/exam-questions/:question_id", handler.GetExamQuestionByID)
	}

	quizPrivate := r.Group("/api")
	quizPrivate.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	{
		// Question Management (Nested under Quiz)
		quizPrivate.POST("/exams/:exam_id/questions", handler.CreateExamQuestion)
		quizPrivate.PUT("/exam-questions/:question_id", handler.UpdateExamQuestion)
		quizPrivate.DELETE("/exam-questions/:question_id", handler.DeleteExamQuestion)
	}
}

type ExamQuestionRequest struct {
	Type          domain.QuestionType `json:"type" binding:"required"`
	Text          string              `json:"text" binding:"required"`
	Options       datatypes.JSON      `json:"options"`
	CorrectAnswer datatypes.JSON      `json:"correct_answer" binding:"required"`
	Points        int                 `json:"points"`
	Explanation   string              `json:"explanation"`
}

func (h *ExamQuestionHandler) CreateExamQuestion(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	var req ExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	question, err := h.examQuestionUsecase.CreateExamQuestion(uint(quizID), req.Type, req.Text, req.Options, req.CorrectAnswer, req.Points, req.Explanation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": question})
}

func (h *ExamQuestionHandler) GetExamQuestionsByExamID(c *gin.Context) {
	quizID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)
	questions, err := h.examQuestionUsecase.GetExamQuestionsByExamID(uint(quizID), true) // true untuk randomize
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": questions})
}

func (h *ExamQuestionHandler) GetExamQuestionByID(c *gin.Context) {
	qID, _ := strconv.ParseUint(c.Param("question_id"), 10, 32)
	question, err := h.examQuestionUsecase.GetExamQuestionByID(uint(qID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Soal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": question})
}

func (h *ExamQuestionHandler) UpdateExamQuestion(c *gin.Context) {
	idParam := c.Param("question_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req ExamQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	question, err := h.examQuestionUsecase.UpdateExamQuestion(uint(id), req.Type, req.Text, req.Options, req.CorrectAnswer, req.Points, req.Explanation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Soal berhasil diperbarui",
		"data":    question,
	})
}

func (h *ExamQuestionHandler) DeleteExamQuestion(c *gin.Context) {
	idParam := c.Param("question_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.examQuestionUsecase.DeleteExamQuestion(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Soal berhasil dihapus"})
}
