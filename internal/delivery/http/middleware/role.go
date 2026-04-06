package middleware

import (
	"github.com/gin-gonic/gin"
)

// RoleMiddleware adalah middleware untuk memeriksa role pengguna yang diizinkan mengakses endpoint tertentu
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(403, gin.H{"error": "Role tidak ditemukan dalam token"})
			c.Abort()
			return
		}
		userRole := role.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}
		c.JSON(403, gin.H{"error": "Akses ditolak: role tidak diizinkan"})
		c.Abort()
	}
}
