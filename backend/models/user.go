package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email" validate:"required,email"`
	Password  string             `bson:"password" json:"-"` // не возвращаем в JSON
	FullName  string             `bson:"fullName" json:"fullName" validate:"required"`
	Phone     string             `bson:"phone" json:"phone"`
	Role      string             `bson:"role" json:"role"` // "user", "admin", "cinema_manager"
	Wallet    Wallet             `bson:"wallet" json:"wallet"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type Wallet struct {
	Balance  float64 `bson:"balance" json:"balance"`
	Currency string  `bson:"currency" json:"currency"` // "KZT"
}
