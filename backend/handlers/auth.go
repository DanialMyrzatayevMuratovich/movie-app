package handlers

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"cinema-booking/utils"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest - структура запроса регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Phone    string `json:"phone"`
}

// LoginRequest - структура запроса логина
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest - структура для обновления профиля
type UpdateProfileRequest struct {
	FullName string `json:"fullName"`
	Phone    string `json:"phone"`
}

// Register - регистрация нового пользователя
func Register(c *gin.Context) {
	var req RegisterRequest

	// 1. Парсинг JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 2. Валидация
	// Email
	req.Email = utils.SanitizeString(req.Email)
	if !utils.ValidateEmail(req.Email) {
		utils.ErrorResponse(c, 400, "Invalid email format")
		return
	}

	// Password
	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		utils.ErrorResponse(c, 400, msg)
		return
	}

	// FullName
	req.FullName = utils.SanitizeString(req.FullName)
	if len(req.FullName) < 2 {
		utils.ErrorResponse(c, 400, "Full name must be at least 2 characters")
		return
	}

	// Phone (опционально)
	if req.Phone != "" {
		req.Phone = utils.SanitizeString(req.Phone)
		if !utils.ValidatePhone(req.Phone) {
			utils.ErrorResponse(c, 400, "Invalid phone format. Use: +7XXXXXXXXXX")
			return
		}
	}

	// 3. Проверить, существует ли email
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		// Пользователь найден - email уже занят
		utils.ErrorResponse(c, 409, "Email already registered")
		return
	}

	// 4. Хешировать пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to hash password")
		return
	}

	// 5. Создать нового пользователя
	newUser := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     "user", // По умолчанию роль "user"
		Wallet: models.Wallet{
			Balance:  0,
			Currency: "KZT",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 6. Сохранить в MongoDB
	result, err := usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to create user")
		return
	}

	// Получить ID созданного пользователя
	newUser.ID = result.InsertedID.(primitive.ObjectID)

	// 7. Создать JWT токен
	token, err := utils.GenerateToken(newUser)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	// 8. Вернуть ответ (без пароля)
	utils.SuccessResponse(c, 201, gin.H{
		"token": token,
		"user": gin.H{
			"id":       newUser.ID.Hex(),
			"email":    newUser.Email,
			"fullName": newUser.FullName,
			"phone":    newUser.Phone,
			"role":     newUser.Role,
			"wallet":   newUser.Wallet,
		},
	})
}

// Login - вход пользователя
func Login(c *gin.Context) {
	var req LoginRequest

	// 1. Парсинг JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 2. Валидация email
	req.Email = utils.SanitizeString(req.Email)
	if !utils.ValidateEmail(req.Email) {
		utils.ErrorResponse(c, 400, "Invalid email format")
		return
	}

	// 3. Найти пользователя в БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		// Пользователь не найден
		utils.ErrorResponse(c, 401, "Invalid email or password")
		return
	}

	// 4. Проверить пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Пароль неправильный
		utils.ErrorResponse(c, 401, "Invalid email or password")
		return
	}

	// 5. Создать JWT токен
	token, err := utils.GenerateToken(user)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to generate token")
		return
	}

	// 6. Вернуть ответ
	utils.SuccessResponse(c, 200, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID.Hex(),
			"email":    user.Email,
			"fullName": user.FullName,
			"phone":    user.Phone,
			"role":     user.Role,
			"wallet":   user.Wallet,
		},
	})
}

// GetProfile - получить профиль текущего пользователя
func GetProfile(c *gin.Context) {
	// 1. Получить userID из контекста (установлен AuthMiddleware)
	userID, exists := c.Get("userId")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized")
		return
	}

	// 2. Конвертировать в ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	// 3. Найти пользователя в БД
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// 4. Вернуть профиль (без пароля)
	utils.SuccessResponse(c, 200, gin.H{
		"id":        user.ID.Hex(),
		"email":     user.Email,
		"fullName":  user.FullName,
		"phone":     user.Phone,
		"role":      user.Role,
		"wallet":    user.Wallet,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

// UpdateProfile - обновить профиль
func UpdateProfile(c *gin.Context) {
	// 1. Получить userID из контекста
	userID, exists := c.Get("userId")
	if !exists {
		utils.ErrorResponse(c, 401, "Unauthorized")
		return
	}

	// 2. Парсинг JSON
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "Invalid request data: "+err.Error())
		return
	}

	// 3. Валидация
	updateFields := bson.M{}

	if req.FullName != "" {
		req.FullName = utils.SanitizeString(req.FullName)
		if len(req.FullName) < 2 {
			utils.ErrorResponse(c, 400, "Full name must be at least 2 characters")
			return
		}
		updateFields["fullName"] = req.FullName
	}

	if req.Phone != "" {
		req.Phone = utils.SanitizeString(req.Phone)
		if !utils.ValidatePhone(req.Phone) {
			utils.ErrorResponse(c, 400, "Invalid phone format. Use: +7XXXXXXXXXX")
			return
		}
		updateFields["phone"] = req.Phone
	}

	// Если нет полей для обновления
	if len(updateFields) == 0 {
		utils.ErrorResponse(c, 400, "No fields to update")
		return
	}

	// Добавить updatedAt
	updateFields["updatedAt"] = time.Now()

	// 4. Обновить в MongoDB
	objectID, _ := primitive.ObjectIDFromHex(userID.(string))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	// Использовать $set оператор
	update := bson.M{"$set": updateFields}

	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to update profile")
		return
	}

	if result.MatchedCount == 0 {
		utils.ErrorResponse(c, 404, "User not found")
		return
	}

	// 5. Получить обновленного пользователя
	var updatedUser models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&updatedUser)
	if err != nil {
		utils.ErrorResponse(c, 500, "Failed to fetch updated profile")
		return
	}

	// 6. Вернуть обновленный профиль
	utils.SuccessWithMessage(c, 200, "Profile updated successfully", gin.H{
		"id":        updatedUser.ID.Hex(),
		"email":     updatedUser.Email,
		"fullName":  updatedUser.FullName,
		"phone":     updatedUser.Phone,
		"role":      updatedUser.Role,
		"wallet":    updatedUser.Wallet,
		"updatedAt": updatedUser.UpdatedAt,
	})
}
