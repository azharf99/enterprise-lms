package middleware

import (
	"github.com/gin-gonic/gin"
)

// RoleMiddleware adalah middleware untuk memeriksa role pengguna yang diizinkan mengakses endpoint tertentu
func RoleMiddleware(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "Role tidak ditemukan dalam token"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "Format role tidak valid"})
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, gin.H{"error": "Akses ditolak: role tidak diizinkan"})
	}
}
