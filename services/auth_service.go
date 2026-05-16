package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/KishiEdward/backend-bengkel.git/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// VerifyFirebaseAndGenerateJWT mengecek keaslian token dari Google, lalu membuat JWT buatan kita
func (s *AuthService) VerifyFirebaseAndGenerateJWT(firebaseToken string) (string, error) {
	// 1. Verifikasi token langsung ke server Google (Firebase)
	// Jika berhasil melewati baris ini, berarti Firebase SDK kamu 100% bekerja!
	token, err := config.FirebaseAuth.VerifyIDToken(context.Background(), firebaseToken)
	if err != nil {
		return "", errors.New("token firebase tidak valid atau sudah kadaluarsa")
	}

	// 2. Jika asli, kita buatkan Token JWT kita sendiri (stempel bengkel)
	claims := jwt.MapClaims{
		"uid": token.UID, // Menyimpan ID user dari Firebase
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 Jam
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	signedToken, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}