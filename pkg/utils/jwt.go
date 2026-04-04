package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPair menyimpan kedua jenis token
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair membuat Access Token (15 menit) dan Refresh Token (7 hari)
func GenerateTokenPair(userID uint, role string) (*TokenPair, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	refreshSecretKey := []byte(os.Getenv("JWT_REFRESH_SECRET")) // Tambahkan ini di .env Anda

	// 1. Buat Access Token (Umur pendek)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	// 2. Buat Refresh Token (Umur panjang)
	// Refresh token biasanya hanya butuh user_id untuk verifikasi
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 Hari
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(refreshSecretKey)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateRefreshToken memeriksa validitas refresh token dan mengembalikan user_id
func ValidateRefreshToken(tokenString string) (uint, error) {
	refreshSecretKey := []byte(os.Getenv("JWT_REFRESH_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshSecretKey, nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("refresh token tidak valid atau kedaluwarsa")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("klaim token tidak valid")
	}

	// Parsing user_id (karena jwt.MapClaims mengubah angka menjadi float64)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("format user_id tidak valid dalam token")
	}

	return uint(userIDFloat), nil
}