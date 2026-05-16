package handlers

import (
	"net/http"

	"github.com/KishiEdward/backend-bengkel.git/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Login menerima token Firebase dari Flutter (atau Postman)
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		FirebaseToken string `json:"firebase_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Firebase token wajib dikirim"})
		return
	}

	// Memanggil otak utama di Auth Service
	jwtToken, err := h.authService.VerifyFirebaseAndGenerateJWT(req.FirebaseToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil, Firebase berfungsi!",
		"token":   jwtToken,
	})
}