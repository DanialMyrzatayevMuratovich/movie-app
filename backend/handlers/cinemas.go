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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetCinemas - получить список кинотеатров с геопоиском
func GetCinemas(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cinemasCollection := config.GetCollection("cinemas")

	// === ПАРАМЕТРЫ ЗАПРОСА ===

	// Пагинация
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}
	skip := (page - 1) * limit

	// Фильтр по городу
	city := c.Query("city") // ?city=Almaty

	// Геопоиск (ближайшие кинотеатры)
	latStr := c.Query("lat")                 // ?lat=43.2220
	lngStr := c.Query("lng")                 // ?lng=76.8512
	maxDistanceStr := c.Query("maxDistance") // ?maxDistance=5000 (метры)

	// === ПОСТРОЕНИЕ ФИЛЬТРА ===

	filter := bson.M{}

	// Фильтр по городу
	if city != "" {
		filter["city"] = city
	}

	// Геопоиск
	var useGeoQuery bool
	var geoQuery bson.M

	if latStr != "" && lngStr != "" {
		lat, errLat := strconv.ParseFloat(latStr, 64)
		lng, errLng := strconv.ParseFloat(lngStr, 64)
		maxDistance := 10000.0 // По умолчанию 10 км

		if maxDistanceStr != "" {
			if dist, err := strconv.ParseFloat(maxDistanceStr, 64); err == nil {
				maxDistance = dist
			}
		}

		if errLat == nil && errLng == nil {
			// Создать геопространственный запрос
			geoQuery = bson.M{
				"location": bson.M{
					"$near": bson.M{
						"$geometry": bson.M{
							"type":        "Point",
							"coordinates": []float64{lng, lat}, // [longitude, latitude]
						},
						"$maxDistance": maxDistance, // в метрах
					},
				},
			}
			useGeoQuery = true
		}
	}

	// === ВЫПОЛНЕНИЕ ЗАПРОСА ===

	var cinemas []models.Cinema
	var total int64

	if useGeoQuery {
		// Геопоиск не поддерживает пагинацию напрямую, используем агрегацию
		pipeline := bson.A{
			bson.M{"$geoNear": bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": geoQuery["location"].(bson.M)["$near"].(bson.M)["$geometry"].(bson.M)["coordinates"],
				},
				"distanceField": "distance",
				"maxDistance":   geoQuery["location"].(bson.M)["$near"].(bson.M)["$maxDistance"],
				"spherical":     true,
			}},
		}

		// Добавить фильтр по городу если есть
		if city != "" {
			pipeline = append(pipeline, bson.M{"$match": bson.M{"city": city}})
		}

		// Пагинация
		pipeline = append(pipeline, bson.M{"$skip": skip})
		pipeline = append(pipeline, bson.M{"$limit": limit})

		cursor, err := cinemasCollection.Aggregate(ctx, pipeline)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to fetch cinemas")
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &cinemas); err != nil {
			utils.ErrorResponse(c, 500, "Failed to decode cinemas")
			return
		}

		// Подсчет общего количества для геопоиска
		countPipeline := bson.A{
			bson.M{"$geoNear": bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": geoQuery["location"].(bson.M)["$near"].(bson.M)["$geometry"].(bson.M)["coordinates"],
				},
				"distanceField": "distance",
				"maxDistance":   geoQuery["location"].(bson.M)["$near"].(bson.M)["$maxDistance"],
				"spherical":     true,
			}},
		}
		if city != "" {
			countPipeline = append(countPipeline, bson.M{"$match": bson.M{"city": city}})
		}
		countPipeline = append(countPipeline, bson.M{"$count": "total"})

		countCursor, _ := cinemasCollection.Aggregate(ctx, countPipeline)
		var countResult []bson.M
		if countCursor != nil {
			countCursor.All(ctx, &countResult)
			if len(countResult) > 0 {
				total = int64(countResult[0]["total"].(int32))
			}
		}

	} else {
		// Обычный запрос без геопоиска
		findOptions := options.Find()
		findOptions.SetSkip(int64(skip))
		findOptions.SetLimit(int64(limit))
		findOptions.SetSort(bson.D{{Key: "name", Value: 1}}) // Сортировка по имени

		cursor, err := cinemasCollection.Find(ctx, filter, findOptions)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to fetch cinemas")
			return
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &cinemas); err != nil {
			utils.ErrorResponse(c, 500, "Failed to decode cinemas")
			return
		}

		// Подсчет общего количества
		total, err = cinemasCollection.CountDocuments(ctx, filter)
		if err != nil {
			utils.ErrorResponse(c, 500, "Failed to count cinemas")
			return
		}
	}

	// === ПОЛУЧИТЬ ИНФОРМАЦИЮ О ЗАЛАХ ===

	// Для каждого кинотеатра получить количество залов
	hallsCollection := config.GetCollection("halls")
	for i := range cinemas {
		hallCount, _ := hallsCollection.CountDocuments(ctx, bson.M{"cinemaId": cinemas[i].ID})
		// Можно добавить дополнительное поле (или использовать len(hallIds))
		_ = hallCount
	}

	// === ОТВЕТ ===

	utils.PaginatedResponse(c, cinemas, page, limit, int(total))
}
