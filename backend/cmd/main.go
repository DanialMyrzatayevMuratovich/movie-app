package main

import (
	"cinema-booking/config"
	"cinema-booking/routes"
	"cinema-booking/scripts"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("🎬 Cinema Booking System - Starting...")

	// 1. Загрузить конфигурацию
	config.LoadConfig()

	// 2. Подключиться к MongoDB
	config.ConnectDB()
	defer config.DisconnectDB()

	// 3. Создать индексы
	log.Println("📊 Creating indexes...")
	scripts.CreateIndexes()

	// 4. Заполнить базу данными 
	// ⚠️ Раскомментируй только при первом запуске!
	// log.Println("🌱 Seeding database...")
	// scripts.SeedDatabase()

	// 5. Настроить Gin
	// Режим (можно поставить gin.ReleaseMode для продакшена)
	gin.SetMode(gin.DebugMode)

	// Создать router
	router := gin.Default()

	// 6. Настроить маршруты
	routes.SetupRoutes(router)

	// 7. Запустить сервер
	port := config.AppConfig.Port
	log.Printf("🚀 Server starting on port %s", port)
	log.Printf("📍 API endpoint: http://localhost:%s/api", port)
	log.Printf("🏥 Health check: http://localhost:%s/api/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
