package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil token dari Header 'Authorization'
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Butuh login untuk akses ini"})
			c.Abort()
			return
		}

		// Format header biasanya: Bearer <token>
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// 2. Validasi Token JWT (cek stempel rahasia di .env)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah expired"})
			c.Abort()
			return
		}

		// Jika valid, izinkan lanjut ke fungsi handler berikutnya
		c.Next()
	}
}