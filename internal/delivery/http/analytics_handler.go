package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsUsecase domain.AnalyticsUsecase
}

func NewAnalyticsHandler(r *gin.Engine, au domain.AnalyticsUsecase) {
	handler := &AnalyticsHandler{analyticsUsecase: au}

	// Route khusus untuk analitik (sebaiknya dilindungi AuthMiddleware khusus role Tutor/Admin)
	analyticsMgmt := r.Group("/api/analytics")
	analyticsMgmt.Use(middleware.RequireAuth(), middleware.RoleMiddleware([]string{"Tutor", "Admin"}))
	// analyticsMgmt.Use(AuthMiddleware(), RoleMiddleware("Tutor", "Admin")) // Contoh proteksi
	{
		analyticsMgmt.GET("/exams/:exam_id/summary", handler.GetExamSummary)
		analyticsMgmt.GET("/exams/:exam_id/item-analysis", handler.GetItemAnalysis)
	}
}

func (h *AnalyticsHandler) GetExamSummary(c *gin.Context) {
	examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)

	analytics, err := h.analyticsUsecase.GetExamAnalytics(uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": analytics})
}

func (h *AnalyticsHandler) GetItemAnalysis(c *gin.Context) {
	examID, _ := strconv.ParseUint(c.Param("exam_id"), 10, 32)

	analysis, err := h.analyticsUsecase.GetItemAnalysis(uint(examID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": analysis})
}
