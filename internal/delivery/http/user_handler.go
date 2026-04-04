package http

import (
	"encoding/csv"
	"net/http"

	"github.com/azharf99/enterprise-lms/internal/delivery/http/middleware"
	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler menyambungkan route Gin dengan Usecase
func NewUserHandler(r *gin.Engine, us domain.UserUsecase) {
	handler := &UserHandler{userUsecase: us}
	// Daftarkan route API di sini
	r.POST("/api/users/login", handler.Login)
	protectedUser := r.Group("/")
	protectedUser.Use(middleware.RequireAuth())
	{
		r.POST("/api/users/import", handler.ImportCSV)
		r.POST("/api/users/refresh", handler.RefreshToken)
	}

}

func (h *UserHandler) ImportCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca format CSV"})
		return
	}

	// Teruskan data ke Usecase
	totalInserted, err := h.userUsecase.ImportFromCSV(records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengimpor pengguna",
		"total":   totalInserted,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	// Validasi input JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	// Teruskan ke Usecase
	token, err := h.userUsecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token tidak disertakan"})
		return
	}

	tokenPair, err := h.userUsecase.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token berhasil diperbarui",
		"tokens":  tokenPair,
	})
}
