package handlers

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"cinema-booking/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateBookingRequest - запрос на создание брони
type CreateBookingRequest struct {
	ShowtimeID    string        `json:"showtimeId" binding:"required"`
	Seats         []SeatRequest `json:"seats" binding:"required,min=1"`
	PaymentMethod string        `json:"paymentMethod" binding:"required"` // "wallet", "card", "cash"
}

type SeatRequest struct {
	Row    string `json:"row" binding:"required"`
	Number int    `json:"number" binding:"required"`
}

// CreateBooking - создать бронь (с MongoDB Transaction)
func CreateBooking(c *gin.Context) {
	// 1. Получить userID из контекста
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	// 2. Парсинг запроса
	var req CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 3. Валидация
	if len(req.Seats) > 10 {
		utils.ErrorResponse(c, 400, "Maximum 10 seats per booking")
		return
	}

	if req.PaymentMethod != "wallet" && req.PaymentMethod != "card" && req.PaymentMethod != "cash" {
		utils.ErrorResponse(c, 400, "Invalid payment method. Use: wallet, card, or cash")
		return
	}

	// 4. Конвертировать showtimeID
	showtimeID, err := primitive.ObjectIDFromHex(req.ShowtimeID)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid showtime ID")
		return
	}

	// === НАЧАЛО MONGODB TRANSACTION ===

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Начать сессию
	session, err := config.MongoClient.StartSession()
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to start database session")
		return
	}
	defer session.EndSession(ctx)

	// Переменные для результата
	var createdBooking models.Booking

	// Выполнить транзакцию
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// === ШАГ 1: Найти сеанс и проверить доступность мест ===

		showtimesCollection := config.GetCollection("showtimes")

		var showtime models.Showtime
		err := showtimesCollection.FindOne(sc, bson.M{"_id": showtimeID}).Decode(&showtime)
		if err != nil {
			return fmt.Errorf("showtime not found")
		}

		// Проверить что сеанс в будущем
		if showtime.StartTime.Before(time.Now()) {
			return fmt.Errorf("showtime has already started")
		}

		// === ШАГ 2: Проверить что места доступны ===

		// Получить информацию о зале
		hallsCollection := config.GetCollection("halls")
		var hall models.Hall
		err = hallsCollection.FindOne(sc, bson.M{"_id": showtime.HallID}).Decode(&hall)
		if err != nil {
			return fmt.Errorf("hall not found")
		}

		// Создать map забронированных мест для быстрой проверки
		bookedSeatsMap := make(map[string]bool)
		for _, seat := range showtime.BookedSeats {
			key := fmt.Sprintf("%s-%d", seat.Row, seat.Number)
			bookedSeatsMap[key] = true
		}

		// Проверить каждое запрошенное место
		var totalAmount float64
		var bookingSeats []models.BookingSeat

		for _, seatReq := range req.Seats {
			key := fmt.Sprintf("%s-%d", seatReq.Row, seatReq.Number)

			// Проверить что место не занято
			if bookedSeatsMap[key] {
				return fmt.Errorf("seat %s-%d is already booked", seatReq.Row, seatReq.Number)
			}

			// Найти место в зале и получить цену
			var seatPrice float64
			seatFound := false

			for _, hallSeat := range hall.Seats {
				if hallSeat.Row == seatReq.Row && hallSeat.Number == seatReq.Number {
					seatPrice = hallSeat.Price + showtime.BasePrice // Цена места + базовая цена сеанса
					seatFound = true
					break
				}
			}

			if !seatFound {
				return fmt.Errorf("seat %s-%d does not exist in this hall", seatReq.Row, seatReq.Number)
			}

			// Добавить место в бронь
			bookingSeats = append(bookingSeats, models.BookingSeat{
				Row:    seatReq.Row,
				Number: seatReq.Number,
				Price:  seatPrice,
			})

			totalAmount += seatPrice
		}

		// === ШАГ 3: Проверить баланс (если оплата через wallet) ===

		usersCollection := config.GetCollection("users")
		var user models.User

		if req.PaymentMethod == "wallet" {
			err = usersCollection.FindOne(sc, bson.M{"_id": userObjectID}).Decode(&user)
			if err != nil {
				return fmt.Errorf("user not found")
			}

			if user.Wallet.Balance < totalAmount {
				return fmt.Errorf("insufficient wallet balance. Required: %.2f KZT, Available: %.2f KZT",
					totalAmount, user.Wallet.Balance)
			}
		}

		// === ШАГ 4: Создать бронь ===

		bookingNumber := fmt.Sprintf("BK-%s-%06d",
			time.Now().Format("20060102"),
			time.Now().UnixNano()%1000000)

		newBooking := models.Booking{
			BookingNumber: bookingNumber,
			UserID:        userObjectID,
			ShowtimeID:    showtimeID,
			Seats:         bookingSeats,
			TotalAmount:   totalAmount,
			Status:        "pending", // Изначально pending
			Payment: models.Payment{
				Method: req.PaymentMethod,
				Status: "pending",
			},
			QRCode:    fmt.Sprintf("QR-%s", bookingNumber),
			ExpiresAt: time.Now().Add(15 * time.Minute), // Истекает через 15 минут
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Если оплата через wallet - сразу подтверждаем
		if req.PaymentMethod == "wallet" {
			newBooking.Status = "confirmed"
			newBooking.Payment.Status = "completed"
			newBooking.Payment.PaidAt = time.Now()
			newBooking.Payment.TransactionID = fmt.Sprintf("TXN-%s", time.Now().Format("20060102150405"))
		}

		bookingsCollection := config.GetCollection("bookings")
		result, err := bookingsCollection.InsertOne(sc, newBooking)
		if err != nil {
			return fmt.Errorf("failed to create booking")
		}

		newBooking.ID = result.InsertedID.(primitive.ObjectID)
		createdBooking = newBooking

		// === ШАГ 5: Обновить сеанс (добавить забронированные места) ===

		// Использовать $push для добавления мест в bookedSeats
		// Использовать $inc для уменьшения availableSeats

		updateSeats := bson.A{}
		for _, seat := range req.Seats {
			updateSeats = append(updateSeats, models.BookedSeat{
				Row:    seat.Row,
				Number: seat.Number,
				Status: "booked",
			})
		}

		_, err = showtimesCollection.UpdateOne(
			sc,
			bson.M{"_id": showtimeID},
			bson.M{
				"$push": bson.M{
					"bookedSeats": bson.M{"$each": updateSeats}, // $push для добавления в массив
				},
				"$inc": bson.M{
					"availableSeats": -len(req.Seats), // $inc для уменьшения счетчика
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update showtime")
		}

		// === ШАГ 6: Списать с кошелька (если wallet) ===

		if req.PaymentMethod == "wallet" {
			// Использовать $inc для уменьшения баланса
			_, err = usersCollection.UpdateOne(
				sc,
				bson.M{"_id": userObjectID},
				bson.M{
					"$inc": bson.M{
						"wallet.balance": -totalAmount, // $inc с отрицательным значением
					},
				},
			)
			if err != nil {
				return fmt.Errorf("failed to deduct from wallet")
			}

			// Создать запись транзакции
			transactionsCollection := config.GetCollection("transactions")
			transaction := models.Transaction{
				UserID:      userObjectID,
				Type:        "booking",
				Amount:      -totalAmount, // Отрицательная сумма (списание)
				BookingID:   newBooking.ID,
				Status:      "completed",
				Description: fmt.Sprintf("Booking payment for %s", bookingNumber),
				CreatedAt:   time.Now(),
			}

			_, err = transactionsCollection.InsertOne(sc, transaction)
			if err != nil {
				return fmt.Errorf("failed to create transaction")
			}
		}

		// Коммит транзакции
		return session.CommitTransaction(sc)
	})

	// === КОНЕЦ TRANSACTION ===

	if err != nil {
		// Транзакция провалилась - все откатилось автоматически
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	// Успешно создано
	utils.SuccessWithMessage(c, 201, "Booking created successfully", createdBooking)
}

// GetMyBookings - получить мои брони
func GetMyBookings(c *gin.Context) {
	// Получить userID из контекста
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	// Фильтр по статусу
	status := c.Query("status") // ?status=confirmed

	filter := bson.M{"userId": userObjectID}
	if status != "" {
		filter["status"] = status
	}

	bookingsCollection := config.GetCollection("bookings")

	// Опции запроса
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Новые первыми

	cursor, err := bookingsCollection.Find(ctx, filter, findOptions)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch bookings")
		return
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		utils.ErrorResponse(c, 500, "Failed to decode bookings")
		return
	}

	// Подсчет общего количества
	total, _ := bookingsCollection.CountDocuments(ctx, filter)

	// Получить детали сеансов (агрегация)
	// Опционально: можно добавить lookup для получения информации о фильме и кинотеатре

	utils.PaginatedResponse(c, bookings, page, limit, int(total))
}

// ConfirmBooking - подтвердить оплату брони (с Transaction)
func ConfirmBooking(c *gin.Context) {
	// Получить userID из контекста
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	// Получить booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid booking ID")
		return
	}

	// === НАЧАЛО TRANSACTION ===

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session, err := config.MongoClient.StartSession()
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to start database session")
		return
	}
	defer session.EndSession(ctx)

	var confirmedBooking models.Booking

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		bookingsCollection := config.GetCollection("bookings")

		// Найти бронь
		var booking models.Booking
		err := bookingsCollection.FindOne(sc, bson.M{
			"_id":    bookingID,
			"userId": userObjectID, // Проверить что бронь принадлежит пользователю
		}).Decode(&booking)

		if err != nil {
			return fmt.Errorf("booking not found")
		}

		// Проверить статус
		if booking.Status == "confirmed" {
			return fmt.Errorf("booking is already confirmed")
		}

		if booking.Status == "cancelled" {
			return fmt.Errorf("booking is cancelled")
		}

		// Проверить что не истекла
		if time.Now().After(booking.ExpiresAt) {
			return fmt.Errorf("booking has expired")
		}

		// Обновить статус брони (использовать $set)
		_, err = bookingsCollection.UpdateOne(
			sc,
			bson.M{"_id": bookingID},
			bson.M{
				"$set": bson.M{
					"status":                "confirmed",
					"payment.status":        "completed",
					"payment.paidAt":        time.Now(),
					"payment.transactionId": fmt.Sprintf("TXN-%s", time.Now().Format("20060102150405")),
					"updatedAt":             time.Now(),
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to confirm booking")
		}

		// Если оплата была не через wallet, нужно создать транзакцию
		if booking.Payment.Method != "wallet" {
			transactionsCollection := config.GetCollection("transactions")
			transaction := models.Transaction{
				UserID:      userObjectID,
				Type:        "booking",
				Amount:      -booking.TotalAmount,
				BookingID:   bookingID,
				Status:      "completed",
				Description: fmt.Sprintf("Booking payment for %s", booking.BookingNumber),
				CreatedAt:   time.Now(),
			}

			_, err = transactionsCollection.InsertOne(sc, transaction)
			if err != nil {
				return fmt.Errorf("failed to create transaction")
			}
		}

		// Получить обновленную бронь
		bookingsCollection.FindOne(sc, bson.M{"_id": bookingID}).Decode(&confirmedBooking)

		return session.CommitTransaction(sc)
	})

	// === КОНЕЦ TRANSACTION ===

	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, 200, "Booking confirmed successfully", confirmedBooking)
}

// CancelBooking - отменить бронь (с Transaction и возвратом денег)
func CancelBooking(c *gin.Context) {
	// Получить userID из контекста
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	// Получить booking ID
	bookingIDStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid booking ID")
		return
	}

	// === НАЧАЛО TRANSACTION ===

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	session, err := config.MongoClient.StartSession()
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to start database session")
		return
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		bookingsCollection := config.GetCollection("bookings")

		// Найти бронь
		var booking models.Booking
		err := bookingsCollection.FindOne(sc, bson.M{
			"_id":    bookingID,
			"userId": userObjectID,
		}).Decode(&booking)

		if err != nil {
			return fmt.Errorf("booking not found")
		}

		// Проверить что можно отменить
		if booking.Status == "cancelled" {
			return fmt.Errorf("booking is already cancelled")
		}

		// Проверить что сеанс еще не начался
		showtimesCollection := config.GetCollection("showtimes")
		var showtime models.Showtime
		err = showtimesCollection.FindOne(sc, bson.M{"_id": booking.ShowtimeID}).Decode(&showtime)
		if err != nil {
			return fmt.Errorf("showtime not found")
		}

		// Нельзя отменить за 2 часа до начала
		if time.Now().Add(2 * time.Hour).After(showtime.StartTime) {
			return fmt.Errorf("cannot cancel booking less than 2 hours before showtime")
		}

		// === ШАГ 1: Обновить статус брони ===

		_, err = bookingsCollection.UpdateOne(
			sc,
			bson.M{"_id": bookingID},
			bson.M{
				"$set": bson.M{
					"status":    "cancelled",
					"updatedAt": time.Now(),
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to cancel booking")
		}

		// === ШАГ 2: Освободить места в сеансе ===

		// Использовать $pull для удаления мест из массива bookedSeats
		// Использовать $inc для увеличения availableSeats

		pullConditions := bson.A{}
		for _, seat := range booking.Seats {
			pullConditions = append(pullConditions, bson.M{
				"row":    seat.Row,
				"number": seat.Number,
			})
		}

		_, err = showtimesCollection.UpdateOne(
			sc,
			bson.M{"_id": booking.ShowtimeID},
			bson.M{
				"$pull": bson.M{
					"bookedSeats": bson.M{
						"$or": pullConditions, // $pull для удаления из массива
					},
				},
				"$inc": bson.M{
					"availableSeats": len(booking.Seats), // $inc для увеличения
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to release seats")
		}

		// === ШАГ 3: Вернуть деньги (если оплачено) ===

		if booking.Status == "confirmed" && booking.Payment.Status == "completed" {
			usersCollection := config.GetCollection("users")

			// Вернуть деньги на кошелек (использовать $inc)
			_, err = usersCollection.UpdateOne(
				sc,
				bson.M{"_id": userObjectID},
				bson.M{
					"$inc": bson.M{
						"wallet.balance": booking.TotalAmount, // $inc с положительным значением
					},
				},
			)
			if err != nil {
				return fmt.Errorf("failed to refund to wallet")
			}

			// Создать запись транзакции (refund)
			transactionsCollection := config.GetCollection("transactions")
			transaction := models.Transaction{
				UserID:      userObjectID,
				Type:        "refund",
				Amount:      booking.TotalAmount, // Положительная сумма (возврат)
				BookingID:   bookingID,
				Status:      "completed",
				Description: fmt.Sprintf("Refund for cancelled booking %s", booking.BookingNumber),
				CreatedAt:   time.Now(),
			}

			_, err = transactionsCollection.InsertOne(sc, transaction)
			if err != nil {
				return fmt.Errorf("failed to create refund transaction")
			}
		}

		return session.CommitTransaction(sc)
	})

	// === КОНЕЦ TRANSACTION ===

	if err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	utils.SuccessWithMessage(c, 200, "Booking cancelled successfully. Refund processed.", gin.H{
		"bookingId": bookingIDStr,
		"cancelled": true,
	})
}
