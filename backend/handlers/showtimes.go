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

// GetShowtimes - получить расписание сеансов с фильтрами
func GetShowtimes(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	showtimesCollection := config.GetCollection("showtimes")

	// === ПАРАМЕТРЫ ЗАПРОСА ===

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	skip := (page - 1) * limit

	// Фильтры
	movieIDStr := c.Query("movieId")   // ?movieId=679ef5a2b3c4d5e6f7890abc
	cinemaIDStr := c.Query("cinemaId") // ?cinemaId=679ef5a2b3c4d5e6f7890def
	hallIDStr := c.Query("hallId")     // ?hallId=679ef5a2b3c4d5e6f7890ghi
	format := c.Query("format")        // ?format=IMAX
	language := c.Query("language")    // ?language=Russian
	dateStr := c.Query("date")         // ?date=2026-02-05

	// Временной диапазон
	fromStr := c.Query("from") // ?from=2026-02-01T10:00:00Z
	toStr := c.Query("to")     // ?to=2026-02-01T23:00:00Z

	// Только будущие сеансы (по умолчанию)
	onlyFuture := c.DefaultQuery("onlyFuture", "true") // ?onlyFuture=false

	// === ПОСТРОЕНИЕ ФИЛЬТРА ===

	filter := bson.M{}

	// Фильтр по фильму
	if movieIDStr != "" {
		movieID, err := primitive.ObjectIDFromHex(movieIDStr)
		if err == nil {
			filter["movieId"] = movieID
		}
	}

	// Фильтр по кинотеатру
	if cinemaIDStr != "" {
		cinemaID, err := primitive.ObjectIDFromHex(cinemaIDStr)
		if err == nil {
			filter["cinemaId"] = cinemaID
		}
	}

	// Фильтр по залу
	if hallIDStr != "" {
		hallID, err := primitive.ObjectIDFromHex(hallIDStr)
		if err == nil {
			filter["hallId"] = hallID
		}
	}

	// Фильтр по формату
	if format != "" {
		filter["format"] = format
	}

	// Фильтр по языку
	if language != "" {
		filter["language"] = language
	}

	// Фильтр по дате (конкретный день)
	if dateStr != "" {
		// Парсинг даты (формат: 2026-02-05)
		date, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			// Начало и конец дня
			startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
			endOfDay := startOfDay.AddDate(0, 0, 1)

			filter["startTime"] = bson.M{
				"$gte": startOfDay,
				"$lt":  endOfDay,
			}
		}
	}

	// Временной диапазон (если date не указан)
	if dateStr == "" {
		timeFilter := bson.M{}

		// От (from)
		if fromStr != "" {
			fromTime, err := time.Parse(time.RFC3339, fromStr)
			if err == nil {
				timeFilter["$gte"] = fromTime
			}
		}

		// До (to)
		if toStr != "" {
			toTime, err := time.Parse(time.RFC3339, toStr)
			if err == nil {
				timeFilter["$lte"] = toTime
			}
		}

		// Только будущие сеансы
		if onlyFuture == "true" && fromStr == "" {
			timeFilter["$gte"] = time.Now()
		}

		// Применить временной фильтр
		if len(timeFilter) > 0 {
			filter["startTime"] = timeFilter
		}
	}

	// === СОРТИРОВКА ===

	sortOptions := bson.D{{Key: "startTime", Value: 1}} // По времени начала (ascending)

	// === ОПЦИИ ЗАПРОСА ===

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(sortOptions)

	// === ВЫПОЛНЕНИЕ ЗАПРОСА ===

	cursor, err := showtimesCollection.Find(ctx, filter, findOptions)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch showtimes")
		return
	}
	defer cursor.Close(ctx)

	var showtimes []models.Showtime
	if err = cursor.All(ctx, &showtimes); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode showtimes")
		return
	}

	// === ПОДСЧЕТ ОБЩЕГО КОЛИЧЕСТВА ===

	total, err := showtimesCollection.CountDocuments(ctx, filter)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to count showtimes")
		return
	}

	// === ДОПОЛНИТЕЛЬНАЯ ИНФОРМАЦИЯ ===

	// Получить информацию о фильмах и кинотеатрах (опционально)
	// Используем агрегацию для объединения данных
	type ShowtimeWithDetails struct {
		models.Showtime `bson:",inline"`
		Movie           models.Movie  `bson:"movie"`
		Cinema          models.Cinema `bson:"cinema"`
	}

	// Если нужна детальная информация
	includeDetails := c.Query("includeDetails") // ?includeDetails=true

	if includeDetails == "true" {
		// Агрегация с lookup
		pipeline := bson.A{
			bson.M{"$match": filter},
			bson.M{"$sort": bson.M{"startTime": 1}},
			bson.M{"$skip": skip},
			bson.M{"$limit": limit},
			// Lookup для фильма
			bson.M{"$lookup": bson.M{
				"from":         "movies",
				"localField":   "movieId",
				"foreignField": "_id",
				"as":           "movieDetails",
			}},
			bson.M{"$unwind": bson.M{
				"path":                       "$movieDetails",
				"preserveNullAndEmptyArrays": true,
			}},
			// Lookup для кинотеатра
			bson.M{"$lookup": bson.M{
				"from":         "cinemas",
				"localField":   "cinemaId",
				"foreignField": "_id",
				"as":           "cinemaDetails",
			}},
			bson.M{"$unwind": bson.M{
				"path":                       "$cinemaDetails",
				"preserveNullAndEmptyArrays": true,
			}},
		}

		aggCursor, err := showtimesCollection.Aggregate(ctx, pipeline)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to fetch detailed showtimes")
			return
		}
		defer aggCursor.Close(ctx)

		var detailedShowtimes []bson.M
		if err = aggCursor.All(ctx, &detailedShowtimes); err != nil {
			utils.ErrorResponse(c, 500, "Failed to decode detailed showtimes")
			return
		}

		utils.PaginatedResponse(c, detailedShowtimes, page, limit, int(total))
		return
	}

	// === ОБЫЧНЫЙ ОТВЕТ ===

	utils.PaginatedResponse(c, showtimes, page, limit, int(total))
}

// CreateShowtime - создать новый сеанс (admin only)
func CreateShowtime(c *gin.Context) {
	var showtime models.Showtime

	// Парсинг JSON
	if err := c.ShouldBindJSON(&showtime); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// Валидация
	if showtime.MovieID.IsZero() {
		utils.ErrorResponse(c, 400, "Movie ID is required")
		return
	}

	if showtime.CinemaID.IsZero() {
		utils.ErrorResponse(c, 400, "Cinema ID is required")
		return
	}

	if showtime.HallID.IsZero() {
		utils.ErrorResponse(c, 400, "Hall ID is required")
		return
	}

	if showtime.StartTime.IsZero() {
		utils.ErrorResponse(c, 400, "Start time is required")
		return
	}

	// Проверить что сеанс в будущем
	if showtime.StartTime.Before(time.Now()) {
		utils.ErrorResponse(c, 400, "Start time must be in the future")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Проверить существование фильма, кинотеатра и зала
	moviesCollection := config.GetCollection("movies")
	var movie models.Movie
	err := moviesCollection.FindOne(ctx, bson.M{"_id": showtime.MovieID}).Decode(&movie)
	if err != nil {
		utils.ErrorResponse(c, 404, "Movie not found")
		return
	}

	cinemasCollection := config.GetCollection("cinemas")
	var cinema models.Cinema
	err = cinemasCollection.FindOne(ctx, bson.M{"_id": showtime.CinemaID}).Decode(&cinema)
	if err != nil {
		utils.ErrorResponse(c, 404, "Cinema not found")
		return
	}

	hallsCollection := config.GetCollection("halls")
	var hall models.Hall
	err = hallsCollection.FindOne(ctx, bson.M{"_id": showtime.HallID, "cinemaId": showtime.CinemaID}).Decode(&hall)
	if err != nil {
		utils.ErrorResponse(c, 404, "Hall not found or doesn't belong to this cinema")
		return
	}

	// Установить значения по умолчанию
	if showtime.EndTime.IsZero() {
		// Рассчитать endTime на основе длительности фильма
		showtime.EndTime = showtime.StartTime.Add(time.Duration(movie.Duration) * time.Minute)
	}

	if showtime.BasePrice == 0 {
		showtime.BasePrice = 2000 // Базовая цена по умолчанию
	}

	showtime.AvailableSeats = hall.Capacity
	showtime.BookedSeats = []models.BookedSeat{} // Пустой массив
	showtime.CreatedAt = time.Now()

	// Сохранить в БД
	showtimesCollection := config.GetCollection("showtimes")

	result, err := showtimesCollection.InsertOne(ctx, showtime)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create showtime")
		return
	}

	showtime.ID = result.InsertedID.(primitive.ObjectID)

	utils.SuccessWithMessage(c, 201, "Showtime created successfully", showtime)
}

// DeleteShowtime - удалить сеанс (admin only)
func DeleteShowtime(c *gin.Context) {
	// Получить ID из URL
	showtimeIDStr := c.Param("id")
	showtimeID, err := primitive.ObjectIDFromHex(showtimeIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid showtime ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	showtimesCollection := config.GetCollection("showtimes")

	// Проверить есть ли брони на этот сеанс
	bookingsCollection := config.GetCollection("bookings")
	bookingCount, _ := bookingsCollection.CountDocuments(ctx, bson.M{
		"showtimeId": showtimeID,
		"status":     bson.M{"$in": []string{"confirmed", "pending"}},
	})

	if bookingCount > 0 {
		utils.ErrorResponse(c, 400, "Cannot delete showtime with existing bookings")
		return
	}

	// Удалить сеанс
	result, err := showtimesCollection.DeleteOne(ctx, bson.M{"_id": showtimeID})
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to delete showtime")
		return
	}

	if result.DeletedCount == 0 {
		utils.ErrorResponse(c, 404, "Showtime not found")
		return
	}

	utils.SuccessWithMessage(c, 200, "Showtime deleted successfully", gin.H{
		"showtimeId": showtimeIDStr,
		"deleted":    true,
	})
}
