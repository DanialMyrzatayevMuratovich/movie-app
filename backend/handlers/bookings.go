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

// CreateBooking - создать бронь
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// === ШАГ 1: Найти сеанс и проверить доступность мест ===

	showtimesCollection := config.GetCollection("showtimes")

	var showtime models.Showtime
	err = showtimesCollection.FindOne(ctx, bson.M{"_id": showtimeID}).Decode(&showtime)
	if err != nil {
		utils.ErrorResponse(c, 404, "Showtime not found")
		return
	}

	// Проверить что сеанс в будущем
	if showtime.StartTime.Before(time.Now()) {
		utils.ErrorResponse(c, 400, "Showtime has already started")
		return
	}

	// === ШАГ 2: Проверить что места доступны ===

	// Попытаться получить информацию о зале
	hallsCollection := config.GetCollection("halls")
	var hall models.Hall
	hallFound := true
	err = hallsCollection.FindOne(ctx, bson.M{"_id": showtime.HallID}).Decode(&hall)
	if err != nil {
		// Зал не найден в БД - используем базовую цену сеанса
		hallFound = false
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
			utils.ErrorResponse(c, 400, fmt.Sprintf("Seat %s-%d is already booked", seatReq.Row, seatReq.Number))
			return
		}

		// Определить цену места
		var seatPrice float64

		if hallFound {
			// Найти место в зале и получить цену
			seatFoundInHall := false
			for _, hallSeat := range hall.Seats {
				if hallSeat.Row == seatReq.Row && hallSeat.Number == seatReq.Number {
					seatPrice = hallSeat.Price + showtime.BasePrice
					seatFoundInHall = true
					break
				}
			}
			if !seatFoundInHall {
				// Место не найдено в зале - используем базовую цену
				seatPrice = showtime.BasePrice
			}
		} else {
			// Зал не найден - используем базовую цену
			seatPrice = showtime.BasePrice
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

	if req.PaymentMethod == "wallet" {
		var user models.User
		err = usersCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&user)
		if err != nil {
			utils.ErrorResponse(c, 404, "User not found")
			return
		}

		if user.Wallet.Balance < totalAmount {
			utils.ErrorResponse(c, 400, fmt.Sprintf("Insufficient wallet balance. Required: %.2f KZT, Available: %.2f KZT",
				totalAmount, user.Wallet.Balance))
			return
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
		Status:        "pending",
		Payment: models.Payment{
			Method: req.PaymentMethod,
			Status: "pending",
		},
		QRCode:    fmt.Sprintf("QR-%s", bookingNumber),
		ExpiresAt: time.Now().Add(15 * time.Minute),
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
	result, err := bookingsCollection.InsertOne(ctx, newBooking)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create booking")
		return
	}

	newBooking.ID = result.InsertedID.(primitive.ObjectID)

	// === ШАГ 5: Обновить сеанс (добавить забронированные места) ===

	updateSeats := bson.A{}
	for _, seat := range req.Seats {
		updateSeats = append(updateSeats, models.BookedSeat{
			Row:    seat.Row,
			Number: seat.Number,
			Status: "booked",
		})
	}

	_, err = showtimesCollection.UpdateOne(
		ctx,
		bson.M{"_id": showtimeID},
		bson.M{
			"$push": bson.M{
				"bookedSeats": bson.M{"$each": updateSeats},
			},
			"$inc": bson.M{
				"availableSeats": -len(req.Seats),
			},
		},
	)
	if err != nil {
		fmt.Printf("Warning: failed to update showtime seats: %v\n", err)
	}

	// === ШАГ 6: Списать с кошелька (если wallet) ===

	if req.PaymentMethod == "wallet" {
		_, err = usersCollection.UpdateOne(
			ctx,
			bson.M{"_id": userObjectID},
			bson.M{
				"$inc": bson.M{
					"wallet.balance": -totalAmount,
				},
			},
		)
		if err != nil {
			fmt.Printf("Warning: failed to deduct from wallet: %v\n", err)
		}

		// Создать запись транзакции
		transactionsCollection := config.GetCollection("transactions")
		transaction := models.Transaction{
			UserID:      userObjectID,
			Type:        "booking",
			Amount:      -totalAmount,
			BookingID:   newBooking.ID,
			Status:      "completed",
			Description: fmt.Sprintf("Booking payment for %s", bookingNumber),
			CreatedAt:   time.Now(),
		}

		_, err = transactionsCollection.InsertOne(ctx, transaction)
		if err != nil {
			fmt.Printf("Warning: failed to create transaction record: %v\n", err)
		}
	}

	// Успешно создано
	utils.SuccessWithMessage(c, 201, "Booking created successfully", newBooking)
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
	status := c.Query("status")

	filter := bson.M{"userId": userObjectID}
	if status != "" {
		filter["status"] = status
	}

	bookingsCollection := config.GetCollection("bookings")

	// Опции запроса
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

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

	utils.PaginatedResponse(c, bookings, page, limit, int(total))
}

// ConfirmBooking - подтвердить оплату брони
func ConfirmBooking(c *gin.Context) {
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	bookingIDStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid booking ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bookingsCollection := config.GetCollection("bookings")

	// Найти бронь
	var booking models.Booking
	err = bookingsCollection.FindOne(ctx, bson.M{
		"_id":    bookingID,
		"userId": userObjectID,
	}).Decode(&booking)

	if err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	if booking.Status == "confirmed" {
		utils.ErrorResponse(c, 400, "Booking is already confirmed")
		return
	}

	if booking.Status == "cancelled" {
		utils.ErrorResponse(c, 400, "Booking is cancelled")
		return
	}

	if time.Now().After(booking.ExpiresAt) {
		utils.ErrorResponse(c, 400, "Booking has expired")
		return
	}

	// Обновить статус брони
	_, err = bookingsCollection.UpdateOne(
		ctx,
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
		utils.ErrorResponse(c, 500, "Failed to confirm booking")
		return
	}

	// Создать транзакцию если оплата не через wallet
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
		transactionsCollection.InsertOne(ctx, transaction)
	}

	// Получить обновленную бронь
	var confirmedBooking models.Booking
	bookingsCollection.FindOne(ctx, bson.M{"_id": bookingID}).Decode(&confirmedBooking)

	utils.SuccessWithMessage(c, 200, "Booking confirmed successfully", confirmedBooking)
}

// CancelBooking - отменить бронь с возвратом денег
func CancelBooking(c *gin.Context) {
	userID, _ := c.Get("userId")
	userObjectID, _ := primitive.ObjectIDFromHex(userID.(string))

	bookingIDStr := c.Param("id")
	bookingID, err := primitive.ObjectIDFromHex(bookingIDStr)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid booking ID")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	bookingsCollection := config.GetCollection("bookings")

	// Найти бронь
	var booking models.Booking
	err = bookingsCollection.FindOne(ctx, bson.M{
		"_id":    bookingID,
		"userId": userObjectID,
	}).Decode(&booking)

	if err != nil {
		utils.ErrorResponse(c, 404, "Booking not found")
		return
	}

	if booking.Status == "cancelled" {
		utils.ErrorResponse(c, 400, "Booking is already cancelled")
		return
	}

	// Проверить что сеанс еще не начался
	showtimesCollection := config.GetCollection("showtimes")
	var showtime models.Showtime
	err = showtimesCollection.FindOne(ctx, bson.M{"_id": booking.ShowtimeID}).Decode(&showtime)
	if err == nil {
		if time.Now().Add(2 * time.Hour).After(showtime.StartTime) {
			utils.ErrorResponse(c, 400, "Cannot cancel booking less than 2 hours before showtime")
			return
		}
	}

	// ШАГ 1: Обновить статус брони
	_, err = bookingsCollection.UpdateOne(
		ctx,
		bson.M{"_id": bookingID},
		bson.M{
			"$set": bson.M{
				"status":    "cancelled",
				"updatedAt": time.Now(),
			},
		},
	)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to cancel booking")
		return
	}

	// ШАГ 2: Освободить места в сеансе
	for _, seat := range booking.Seats {
		showtimesCollection.UpdateOne(
			ctx,
			bson.M{"_id": booking.ShowtimeID},
			bson.M{
				"$pull": bson.M{
					"bookedSeats": bson.M{
						"row":    seat.Row,
						"number": seat.Number,
					},
				},
				"$inc": bson.M{
					"availableSeats": 1,
				},
			},
		)
	}

	// ШАГ 3: Вернуть деньги (если оплачено)
	if booking.Status == "confirmed" && booking.Payment.Status == "completed" {
		usersCollection := config.GetCollection("users")

		_, err = usersCollection.UpdateOne(
			ctx,
			bson.M{"_id": userObjectID},
			bson.M{
				"$inc": bson.M{
					"wallet.balance": booking.TotalAmount,
				},
			},
		)
		if err != nil {
			fmt.Printf("Warning: failed to refund to wallet: %v\n", err)
		}

		// Создать запись транзакции (refund)
		transactionsCollection := config.GetCollection("transactions")
		transaction := models.Transaction{
			UserID:      userObjectID,
			Type:        "refund",
			Amount:      booking.TotalAmount,
			BookingID:   bookingID,
			Status:      "completed",
			Description: fmt.Sprintf("Refund for cancelled booking %s", booking.BookingNumber),
			CreatedAt:   time.Now(),
		}
		transactionsCollection.InsertOne(ctx, transaction)
	}

	utils.SuccessWithMessage(c, 200, "Booking cancelled successfully. Refund processed.", gin.H{
		"bookingId": bookingIDStr,
		"cancelled": true,
	})
}
