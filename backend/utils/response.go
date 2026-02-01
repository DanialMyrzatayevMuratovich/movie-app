package utils

import "github.com/gin-gonic/gin"

// SuccessResponse - стандартный успешный ответ
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"data":    data,
	})
}

// ErrorResponse - стандартный ответ с ошибкой
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   message,
	})
}

// SuccessWithMessage - успешный ответ с сообщением
func SuccessWithMessage(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// PaginatedResponse - ответ с пагинацией
func PaginatedResponse(c *gin.Context, data interface{}, page, limit, total int) {
	c.JSON(200, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	})
}
