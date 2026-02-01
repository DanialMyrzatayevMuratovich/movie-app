package middleware

import (
	"cinema-booking/utils"
	"log"

	"github.com/gin-gonic/gin"
)

// ErrorHandler - middleware для обработки паник и ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ Panic recovered: %v", err)
				utils.ErrorResponse(c, 500, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()

		// Проверить если были ошибки в обработчике
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("❌ Error: %v", err.Error())

			// Если ответ еще не отправлен
			if !c.Writer.Written() {
				utils.ErrorResponse(c, 500, err.Error())
			}
		}
	}
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
