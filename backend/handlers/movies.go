package handlers

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"cinema-booking/utils"
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetMovies - получить список фильмов с фильтрами, поиском и пагинацией
func GetMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moviesCollection := config.GetCollection("movies")

	// === ПАРАМЕТРЫ ЗАПРОСА ===

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Поиск (full-text search)
	search := c.Query("search") // ?search=Avatar

	// Фильтры
	genre := c.Query("genre")         // ?genre=Action
	language := c.Query("language")   // ?language=English
	isActive := c.Query("isActive")   // ?isActive=true
	minRating := c.Query("minRating") // ?minRating=8.0

	// Сортировка
	sortBy := c.DefaultQuery("sortBy", "createdAt")  // ?sortBy=imdbRating
	sortOrder := c.DefaultQuery("sortOrder", "desc") // ?sortOrder=asc

	// === ПОСТРОЕНИЕ ФИЛЬТРА ===

	filter := bson.M{}

	// Только активные фильмы (по умолчанию)
	switch isActive {
	case "", "true":
		filter["isActive"] = true
	case "false":
		filter["isActive"] = false
	}

	// Полнотекстовый поиск
	if search != "" {
		filter["$text"] = bson.M{"$search": search}
	}

	// Фильтр по жанру
	if genre != "" {
		filter["genres"] = genre // MongoDB автоматически ищет в массиве
	}

	// Фильтр по языку
	if language != "" {
		filter["language"] = language
	}

	// Фильтр по рейтингу
	if minRating != "" {
		rating, err := strconv.ParseFloat(minRating, 64)
		if err == nil {
			filter["imdbRating"] = bson.M{"$gte": rating}
		}
	}

	// === СОРТИРОВКА ===

	sortDirection := 1 // ascending
	if sortOrder == "desc" {
		sortDirection = -1
	}

	sortOptions := bson.D{{Key: sortBy, Value: sortDirection}}

	// === ОПЦИИ ЗАПРОСА ===

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sortOptions)

	// === ВЫПОЛНЕНИЕ ЗАПРОСА ===

	cursor, err := moviesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch movies")
		return
	}
	defer cursor.Close(ctx)

	// Декодировать результаты
	var movies []models.Movie
	if err = cursor.All(ctx, &movies); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode movies")
		return
	}

	// === ПОДСЧЕТ ОБЩЕГО КОЛИЧЕСТВА ===

	total, err := moviesCollection.CountDocuments(ctx, filter)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to count movies")
		return
	}

	// === ОТВЕТ С ПАГИНАЦИЕЙ ===

	utils.PaginatedResponse(c, movies, page, limit, int(total))
}

// GetMovieDetails - получить детальную информацию о фильме
func GetMovieDetails(c *gin.Context) {
	// Получить ID из URL параметра
	movieIDStr := c.Param("id")

	// Конвертировать в ObjectID
	movieID, err := primitive.ObjectIDFromHex(movieIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid movie ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moviesCollection := config.GetCollection("movies")

	// Найти фильм
	var movie models.Movie
	err = moviesCollection.FindOne(ctx, bson.M{"_id": movieID}).Decode(&movie)
	if err != nil {
		utils.ErrorResponse(c, 404, "Movie not found")
		return
	}

	// === ДОПОЛНИТЕЛЬНАЯ ИНФОРМАЦИЯ ===

	// Подсчитать средний рейтинг из отзывов
	var averageReviewRating float64
	if len(movie.Reviews) > 0 {
		var sum float64
		for _, review := range movie.Reviews {
			sum += review.Rating
		}
		averageReviewRating = sum / float64(len(movie.Reviews))
	}

	// Найти ближайшие сеансы (опционально)
	showtimesCollection := config.GetCollection("showtimes")

	// Найти сеансы этого фильма в ближайшие 7 дней
	now := time.Now()
	sevenDaysLater := now.AddDate(0, 0, 7)

	showtimeFilter := bson.M{
		"movieId": movieID,
		"startTime": bson.M{
			"$gte": now,
			"$lte": sevenDaysLater,
		},
	}

	showtimeCursor, err := showtimesCollection.Find(ctx, showtimeFilter, options.Find().SetLimit(5))
	if err == nil {
		var showtimes []bson.M
		showtimeCursor.All(ctx, &showtimes)
		showtimeCursor.Close(ctx)

		// Вернуть ответ с дополнительной информацией
		utils.SuccessResponse(c, 200, gin.H{
			"movie":               movie,
			"averageReviewRating": averageReviewRating,
			"totalReviews":        len(movie.Reviews),
			"upcomingShowtimes":   len(showtimes),
		})
		return
	}

	// Если не удалось получить сеансы, вернуть только фильм
	utils.SuccessResponse(c, 200, gin.H{
		"movie":               movie,
		"averageReviewRating": averageReviewRating,
		"totalReviews":        len(movie.Reviews),
	})
}

// CreateMovie - создать новый фильм (admin only)
func CreateMovie(c *gin.Context) {
	var movie models.Movie

	// Парсинг JSON
	if err := c.ShouldBindJSON(&movie); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// Валидация
	if movie.Title == "" {
		utils.ErrorResponse(c, 400, "Title is required")
		return
	}

	// Установить значения по умолчанию
	movie.IsActive = true
	movie.CreatedAt = time.Now()
	movie.Reviews = []models.Review{} // Пустой массив отзывов

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moviesCollection := config.GetCollection("movies")

	// Сохранить в БД
	result, err := moviesCollection.InsertOne(ctx, movie)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create movie")
		return
	}

	movie.ID = result.InsertedID.(primitive.ObjectID)

	utils.SuccessWithMessage(c, 201, "Movie created successfully", movie)
}

// UpdateMovie - обновить фильм (admin only)
func UpdateMovie(c *gin.Context) {
	// Получить ID из URL
	movieIDStr := c.Param("id")
	movieID, err := primitive.ObjectIDFromHex(movieIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid movie ID")
		return
	}

	// Парсинг данных для обновления
	var updateData bson.M
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// Удалить поля, которые нельзя обновлять напрямую
	delete(updateData, "_id")
	delete(updateData, "createdAt")
	delete(updateData, "reviews") // Отзывы обновляются отдельно

	// Если нет полей для обновления
	if len(updateData) == 0 {
		utils.ErrorResponse(c, 400, "No fields to update")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moviesCollection := config.GetCollection("movies")

	// Использовать $set оператор
	update := bson.M{"$set": updateData}

	result, err := moviesCollection.UpdateOne(ctx, bson.M{"_id": movieID}, update)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update movie")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, 404, "Movie not found")
		return
	}

	// Получить обновленный фильм
	var updatedMovie models.Movie
	moviesCollection.FindOne(ctx, bson.M{"_id": movieID}).Decode(&updatedMovie)

	utils.SuccessWithMessage(c, 200, "Movie updated successfully", updatedMovie)
}

// DeleteMovie - удалить фильм (деактивировать)
func DeleteMovie(c *gin.Context) {
	// Получить ID из URL
	movieIDStr := c.Param("id")
	movieID, err := primitive.ObjectIDFromHex(movieIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid movie ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	moviesCollection := config.GetCollection("movies")

	// Soft delete - просто деактивировать фильм
	update := bson.M{"$set": bson.M{"isActive": false}}

	result, err := moviesCollection.UpdateOne(ctx, bson.M{"_id": movieID}, update)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete movie")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, 404, "Movie not found")
		return
	}

	utils.SuccessWithMessage(c, 200, "Movie deleted successfully", gin.H{
		"movieId": movieIDStr,
		"deleted": true,
	})
}
