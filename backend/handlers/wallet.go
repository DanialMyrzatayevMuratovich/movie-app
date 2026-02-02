package handlers

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"cinema-booking/utils"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TopUpRequest - запрос на пополнение кошелька
type TopUpRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// TopUpWallet - пополнить кошелек пользователя
func TopUpWallet(c *gin.Context) {
	// 1. Получить userID из контекста
	userID, exists := c.Get("userId")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized")
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	// 2. Парсинг запроса
	var req TopUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 3. Валидация суммы
	if req.Amount <= 0 {
		utils.ErrorResponse(c, 400, "Amount must be greater than 0")
		return
	}

	if req.Amount > 1000000 {
		utils.ErrorResponse(c, 400, "Maximum top-up amount is 1,000,000 KZT")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	// 4. Пополнить баланс используя $inc
	result, err := usersCollection.UpdateOne(
		ctx,
		bson.M{"_id": userObjectID},
		bson.M{
			"$inc": bson.M{
				"wallet.balance": req.Amount,
			},
			"$set": bson.M{
				"updatedAt": time.Now(),
			},
		},
	)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to top up wallet")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// 5. Создать запись транзакции
	transactionsCollection := config.GetCollection("transactions")
	transaction := models.Transaction{
		UserID:      userObjectID,
		Type:        "wallet_topup",
		Amount:      req.Amount,
		Status:      "completed",
		Description: fmt.Sprintf("Wallet top-up: +%.2f KZT", req.Amount),
		CreatedAt:   time.Now(),
	}

	_, err = transactionsCollection.InsertOne(ctx, transaction)
	if err != nil {
		// Транзакция записи не критична, продолжаем
		fmt.Printf("Warning: failed to create transaction record: %v\n", err)
	}

	// 6. Получить обновленного пользователя
	var updatedUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&updatedUser)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch updated user")
		return
	}

	utils.SuccessWithMessage(c, 200, "Wallet topped up successfully", gin.H{
		"wallet": updatedUser.Wallet,
		"amount": req.Amount,
	})
}
