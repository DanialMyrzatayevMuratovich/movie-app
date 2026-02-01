package utils

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims - структура для JWT токена
type JWTClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken - создать JWT токен для пользователя
func GenerateToken(user models.User) (string, error) {
	// Получить срок действия из конфига
	expirationTime := time.Now().Add(24 * time.Hour) // 24 часа

	claims := JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "cinema-booking-system",
		},
	}

	// Создать токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписать токен секретным ключом
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken - проверить JWT токен
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверить алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Получить claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
