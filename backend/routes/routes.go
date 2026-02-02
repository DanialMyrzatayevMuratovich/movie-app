package routes

import (
	"cinema-booking/handlers"
	"cinema-booking/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes - настроить все маршруты
func SetupRoutes(router *gin.Engine) {
	// Применить глобальные middleware
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORSMiddleware())

	// API группа
	api := router.Group("/api")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Cinema Booking API is running",
			})
		})

		// Auth routes (публичные)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Movies (публичные - можно смотреть без авторизации)
		api.GET("/movies", handlers.GetMovies)
		api.GET("/movies/:id", handlers.GetMovieDetails)

		// Cinemas (публичные)
		api.GET("/cinemas", handlers.GetCinemas)

		// Showtimes (публичные)
		api.GET("/showtimes", handlers.GetShowtimes)

		// Protected routes (требуют авторизации)
		authorized := api.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// Profile
			authorized.GET("/profile", handlers.GetProfile)
			authorized.PUT("/profile", handlers.UpdateProfile)

			// Wallet
			authorized.POST("/wallet/topup", handlers.TopUpWallet)

			// Bookings (только для авторизованных пользователей)
			authorized.POST("/bookings", handlers.CreateBooking)
			authorized.GET("/bookings/my", handlers.GetMyBookings)
			authorized.POST("/bookings/:id/confirm", handlers.ConfirmBooking)
			authorized.DELETE("/bookings/:id", handlers.CancelBooking)

			// Analytics
			authorized.GET("/analytics/popular-movies", handlers.GetPopularMovies)
			authorized.GET("/analytics/cinema-stats", handlers.GetCinemaStats)
			authorized.GET("/analytics/revenue", handlers.GetRevenue)
		}

		// Admin routes (только для админов)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware())
		admin.Use(middleware.RequireRole("admin"))
		{
			// Управление фильмами
			admin.POST("/movies", handlers.CreateMovie)
			admin.PUT("/movies/:id", handlers.UpdateMovie)
			admin.DELETE("/movies/:id", handlers.DeleteMovie)

			// Управление сеансами
			admin.POST("/showtimes", handlers.CreateShowtime)
			admin.DELETE("/showtimes/:id", handlers.DeleteShowtime)
		}
	}
}
