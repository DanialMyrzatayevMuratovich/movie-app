package middleware

import (
	"cinema-booking/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole - middleware для проверки роли пользователя
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получить роль из контекста (должна быть установлена AuthMiddleware)
		userRole, exists := c.Get("userRole")
		if !exists {
			utils.ErrorResponse(c, 403, "Forbidden: No role information")
			c.Abort()
			return
		}

		roleStr := userRole.(string)

		// Проверить, есть ли роль пользователя в списке разрешенных
		for _, role := range allowedRoles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		// Роль не подходит
		utils.ErrorResponse(c, 403, "Forbidden: Insufficient permissions")
		c.Abort()
	}
}
