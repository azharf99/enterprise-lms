package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders menambahkan header keamanan standar industri
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

// SetupCORS membatasi domain mana saja yang boleh mengakses API ini
func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Nanti ganti dengan domain frontend Anda yang sebenarnya
		AllowOrigins:     []string{"https://lms.azharfa.cloud", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// --- Rate Limiter In-Memory ---
// Melindungi server dari spam request (misal: frontend stuck di infinite loop refresh token)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Mengizinkan 5 request per detik, dengan burst maksimal 10 request
		limiter = rate.NewLimiter(5, 10)
		visitors[ip] = limiter
	}
	return limiter
}

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Terlalu banyak permintaan, silakan coba beberapa saat lagi.",
			})
			return
		}
		c.Next()
	}
}
