package handlers

import (
	"cinema-booking/config"
	"cinema-booking/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetPopularMovies - топ популярных фильмов (Aggregation Pipeline)
func GetPopularMovies(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Параметры
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Временной диапазон (по умолчанию последние 30 дней)
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)
	fromDate := time.Now().AddDate(0, 0, -days)

	bookingsCollection := config.GetCollection("bookings")

	// === АГРЕГАЦИОННЫЙ PIPELINE ===

	pipeline := bson.A{
		// Шаг 1: Фильтр по дате и статусу
		bson.M{"$match": bson.M{
			"status":    "confirmed",
			"createdAt": bson.M{"$gte": fromDate},
		}},

		// Шаг 2: Lookup - получить информацию о сеансе
		bson.M{"$lookup": bson.M{
			"from":         "showtimes",
			"localField":   "showtimeId",
			"foreignField": "_id",
			"as":           "showtime",
		}},

		// Шаг 3: Unwind - развернуть массив showtimes
		bson.M{"$unwind": "$showtime"},

		// Шаг 4: Lookup - получить информацию о фильме
		bson.M{"$lookup": bson.M{
			"from":         "movies",
			"localField":   "showtime.movieId",
			"foreignField": "_id",
			"as":           "movie",
		}},

		// Шаг 5: Unwind - развернуть массив movies
		bson.M{"$unwind": "$movie"},

		// Шаг 6: Group - группировка по фильму
		bson.M{"$group": bson.M{
			"_id":        "$movie._id",
			"title":      bson.M{"$first": "$movie.title"},
			"titleKz":    bson.M{"$first": "$movie.titleKz"},
			"titleRu":    bson.M{"$first": "$movie.titleRu"},
			"poster":     bson.M{"$first": "$movie.posterFileId"},
			"genres":     bson.M{"$first": "$movie.genres"},
			"imdbRating": bson.M{"$first": "$movie.imdbRating"},

			// Агрегированные метрики
			"totalBookings": bson.M{"$sum": 1},
			"totalRevenue":  bson.M{"$sum": "$totalAmount"},
			"totalTickets":  bson.M{"$sum": bson.M{"$size": "$seats"}},

			// Средняя цена билета
			"averageTicketPrice": bson.M{"$avg": bson.M{
				"$divide": bson.A{"$totalAmount", bson.M{"$size": "$seats"}},
			}},

			// Самая ранняя и поздняя бронь
			"firstBooking": bson.M{"$min": "$createdAt"},
			"lastBooking":  bson.M{"$max": "$createdAt"},
		}},

		// Шаг 7: Sort - сортировка по количеству броней
		bson.M{"$sort": bson.M{"totalBookings": -1}},

		// Шаг 8: Limit - ограничить результат
		bson.M{"$limit": limit},

		// Шаг 9: Project - выбрать нужные поля и добавить вычисления
		bson.M{"$project": bson.M{
			"_id":                1,
			"title":              1,
			"titleKz":            1,
			"titleRu":            1,
			"poster":             1,
			"genres":             1,
			"imdbRating":         1,
			"totalBookings":      1,
			"totalRevenue":       1,
			"totalTickets":       1,
			"averageTicketPrice": bson.M{"$round": bson.A{"$averageTicketPrice", 2}},
			"firstBooking":       1,
			"lastBooking":        1,

			// Популярность (баллы)
			"popularityScore": bson.M{"$add": bson.A{
				bson.M{"$multiply": bson.A{"$totalBookings", 10}},
				bson.M{"$multiply": bson.A{"$totalTickets", 5}},
			}},
		}},
	}

	// Выполнить агрегацию
	cursor, err := bookingsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to aggregate popular movies: "+err.Error())
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode results")
		return
	}

	utils.SuccessResponse(c, 200, gin.H{
		"movies":      results,
		"period":      fmt.Sprintf("Last %d days", days),
		"totalMovies": len(results),
	})
}

// GetCinemaStats - статистика по кинотеатрам (Aggregation Pipeline)
func GetCinemaStats(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Параметры
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)
	fromDate := time.Now().AddDate(0, 0, -days)

	bookingsCollection := config.GetCollection("bookings")

	// === АГРЕГАЦИОННЫЙ PIPELINE ===

	null := 0
	pipeline := bson.A{
		// Шаг 1: Фильтр по дате и статусу
		bson.M{"$match": bson.M{
			"status":    "confirmed",
			"createdAt": bson.M{"$gte": fromDate},
		}},

		// Шаг 2: Lookup - получить сеанс
		bson.M{"$lookup": bson.M{
			"from":         "showtimes",
			"localField":   "showtimeId",
			"foreignField": "_id",
			"as":           "showtime",
		}},
		bson.M{"$unwind": "$showtime"},

		// Шаг 3: Lookup - получить кинотеатр
		bson.M{"$lookup": bson.M{
			"from":         "cinemas",
			"localField":   "showtime.cinemaId",
			"foreignField": "_id",
			"as":           "cinema",
		}},
		bson.M{"$unwind": "$cinema"},

		// Шаг 4: Group - группировка по кинотеатру
		bson.M{"$group": bson.M{
			"_id":        "$cinema._id",
			"cinemaName": bson.M{"$first": "$cinema.name"},
			"city":       bson.M{"$first": "$cinema.city"},
			"address":    bson.M{"$first": "$cinema.address"},
			"rating":     bson.M{"$first": "$cinema.rating"},

			// Метрики
			"totalBookings": bson.M{"$sum": 1},
			"totalRevenue":  bson.M{"$sum": "$totalAmount"},
			"totalTickets":  bson.M{"$sum": bson.M{"$size": "$seats"}},

			// Средняя цена билета
			"averageTicketPrice": bson.M{"$avg": bson.M{
				"$divide": bson.A{"$totalAmount", bson.M{"$size": "$seats"}},
			}},

			// Средний размер брони (мест на бронь)
			"averageSeatsPerBooking": bson.M{"$avg": bson.M{"$size": "$seats"}},
		}},

		// Шаг 5: Добавить процент от общей выручки
		bson.M{"$group": bson.M{
			"_id":             null,
			"cinemas":         bson.M{"$push": "$$ROOT"},
			"totalRevenueAll": bson.M{"$sum": "$totalRevenue"},
		}},

		// Шаг 6: Unwind для расчета процентов
		bson.M{"$unwind": "$cinemas"},

		// Шаг 7: Project - добавить процент выручки
		bson.M{"$project": bson.M{
			"_id":                    "$cinemas._id",
			"cinemaName":             "$cinemas.cinemaName",
			"city":                   "$cinemas.city",
			"address":                "$cinemas.address",
			"rating":                 "$cinemas.rating",
			"totalBookings":          "$cinemas.totalBookings",
			"totalRevenue":           "$cinemas.totalRevenue",
			"totalTickets":           "$cinemas.totalTickets",
			"averageTicketPrice":     bson.M{"$round": bson.A{"$cinemas.averageTicketPrice", 2}},
			"averageSeatsPerBooking": bson.M{"$round": bson.A{"$cinemas.averageSeatsPerBooking", 2}},

			// Процент от общей выручки
			"revenuePercentage": bson.M{"$round": bson.A{
				bson.M{"$multiply": bson.A{
					bson.M{"$divide": bson.A{"$cinemas.totalRevenue", "$totalRevenueAll"}},
					100,
				}},
				2,
			}},
		}},

		// Шаг 8: Sort - сортировка по выручке
		bson.M{"$sort": bson.M{"totalRevenue": -1}},
	}

	// Выполнить агрегацию
	cursor, err := bookingsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to aggregate cinema stats: "+err.Error())
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode results")
		return
	}

	// Подсчитать общую статистику
	var totalRevenue, totalTickets, totalBookings float64
	for _, cinema := range results {
		totalRevenue += cinema["totalRevenue"].(float64)
		totalTickets += float64(cinema["totalTickets"].(int32))
		totalBookings += float64(cinema["totalBookings"].(int32))
	}

	utils.SuccessResponse(c, 200, gin.H{
		"cinemas": results,
		"period":  fmt.Sprintf("Last %d days", days),
		"summary": gin.H{
			"totalCinemas":  len(results),
			"totalRevenue":  totalRevenue,
			"totalTickets":  totalTickets,
			"totalBookings": totalBookings,
		},
	})
}

// GetRevenue - статистика выручки по периодам (Aggregation Pipeline)
func GetRevenue(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Параметры
	groupBy := c.DefaultQuery("groupBy", "day") // day, week, month
	daysStr := c.DefaultQuery("days", "30")
	days, _ := strconv.Atoi(daysStr)
	fromDate := time.Now().AddDate(0, 0, -days)

	// Определить формат группировки
	var dateFormat string
	var groupByDate bson.M

	switch groupBy {
	case "day":
		dateFormat = "%Y-%m-%d"
		groupByDate = bson.M{
			"year":  bson.M{"$year": "$createdAt"},
			"month": bson.M{"$month": "$createdAt"},
			"day":   bson.M{"$dayOfMonth": "$createdAt"},
		}
	case "week":
		dateFormat = "%Y-W%V"
		groupByDate = bson.M{
			"year": bson.M{"$year": "$createdAt"},
			"week": bson.M{"$week": "$createdAt"},
		}
	case "month":
		dateFormat = "%Y-%m"
		groupByDate = bson.M{
			"year":  bson.M{"$year": "$createdAt"},
			"month": bson.M{"$month": "$createdAt"},
		}
	default:
		dateFormat = "%Y-%m-%d"
		groupByDate = bson.M{
			"year":  bson.M{"$year": "$createdAt"},
			"month": bson.M{"$month": "$createdAt"},
			"day":   bson.M{"$dayOfMonth": "$createdAt"},
		}
	}

	bookingsCollection := config.GetCollection("bookings")

	// === АГРЕГАЦИОННЫЙ PIPELINE ===

	pipeline := bson.A{
		// Шаг 1: Фильтр по дате и статусу
		bson.M{"$match": bson.M{
			"status":    "confirmed",
			"createdAt": bson.M{"$gte": fromDate},
		}},

		// Шаг 2: Group - группировка по периоду
		bson.M{"$group": bson.M{
			"_id": groupByDate,
			"date": bson.M{"$first": bson.M{
				"$dateToString": bson.M{
					"format": dateFormat,
					"date":   "$createdAt",
				},
			}},

			// Метрики
			"totalRevenue":  bson.M{"$sum": "$totalAmount"},
			"totalBookings": bson.M{"$sum": 1},
			"totalTickets":  bson.M{"$sum": bson.M{"$size": "$seats"}},

			// Средняя стоимость брони
			"averageBookingValue": bson.M{"$avg": "$totalAmount"},

			// Минимальная и максимальная бронь
			"minBookingValue": bson.M{"$min": "$totalAmount"},
			"maxBookingValue": bson.M{"$max": "$totalAmount"},
		}},

		// Шаг 3: Sort - сортировка по дате
		bson.M{"$sort": bson.M{"_id": 1}},

		// Шаг 4: Project - форматирование результата
		bson.M{"$project": bson.M{
			"_id":                 0,
			"date":                "$date",
			"totalRevenue":        bson.M{"$round": bson.A{"$totalRevenue", 2}},
			"totalBookings":       1,
			"totalTickets":        1,
			"averageBookingValue": bson.M{"$round": bson.A{"$averageBookingValue", 2}},
			"minBookingValue":     bson.M{"$round": bson.A{"$minBookingValue", 2}},
			"maxBookingValue":     bson.M{"$round": bson.A{"$maxBookingValue", 2}},
		}},
	}

	// Выполнить агрегацию
	cursor, err := bookingsCollection.Aggregate(ctx, pipeline)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to aggregate revenue: "+err.Error())
		return
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode results")
		return
	}

	// Подсчитать общую статистику
	var totalRevenue, totalTickets, totalBookings float64
	for _, period := range results {
		totalRevenue += period["totalRevenue"].(float64)
		totalTickets += float64(period["totalTickets"].(int32))
		totalBookings += float64(period["totalBookings"].(int32))
	}

	utils.SuccessResponse(c, 200, gin.H{
		"data":    results,
		"groupBy": groupBy,
		"period":  fmt.Sprintf("Last %d days", days),
		"summary": gin.H{
			"totalRevenue":            totalRevenue,
			"totalTickets":            totalTickets,
			"totalBookings":           totalBookings,
			"averageRevenuePerPeriod": totalRevenue / float64(len(results)),
		},
	})
}
