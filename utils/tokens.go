package utils

import (
	"go_edtech_backend/db"
	"go_edtech_backend/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateNewTokens(email string) (string, string, error) {
	var existingUser models.User
	if err := db.DB.Where("email = ?", email).First(&existingUser).Error; err != nil {
		return "", "", err
	}

	token, err := generateJWT(email, time.Now().Add(time.Hour*48)) // Valid for 48 hours
	if err != nil {
		return "", "", err
	}
	refreshToken, err := generateJWT(email, time.Now().Add(time.Hour*240)) // Valid for 10 days
	if err != nil {
		return "", "", err
	}

	if err := db.DB.Model(&existingUser).Updates(models.User{Token: token, RefreshToken: refreshToken}).Error; err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func generateJWT(email string, expiryTime time.Time) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   expiryTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := db.GetJWTSecret()
	return token.SignedString(jwtSecret)
}

func VerifyToken(token string) (*jwt.Token, error) { // expired or not
	key := db.GetJWTSecret()
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	return tk, err
}
