package http

import (
	"encoding/csv"
	"net/http"
	"strconv"

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
	r.POST("/api/users/refresh", handler.RefreshToken)
	protectedUser := r.Group("/api/users")
	protectedUser.Use(middleware.RequireAuth())
	{
		protectedUser.GET("", handler.GetAll)
		protectedUser.POST("", handler.CreateUser)
		protectedUser.PUT("/:user_id", handler.UpdateUser)
		protectedUser.DELETE("/:user_id", handler.DeleteUser)
		protectedUser.POST("/import", handler.ImportCSV)
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

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}
	user, err := h.userUsecase.CreateUser(req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pengguna berhasil dibuat", "data": user})
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

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("user_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var req domain.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format request tidak valid"})
		return
	}

	user, err := h.userUsecase.UpdateUser(uint(id), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengguna berhasil diperbarui",
		"data":    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("user_id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.userUsecase.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
