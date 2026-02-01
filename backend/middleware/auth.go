package middleware

import (
	"cinema-booking/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware для проверки JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Получить токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "Authorization header is required")
			c.Abort()
			return
		}

		// 2. Проверить формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, 401, "Invalid authorization header format. Use: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Валидировать токен
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, 401, "Invalid or expired token")
			c.Abort()
			return
		}

		// 4. Сохранить данные пользователя в контекст
		c.Set("userId", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		// Продолжить выполнение
		c.Next()
	}
}
