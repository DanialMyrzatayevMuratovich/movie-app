package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"userId" json:"userId"`
	Type        string             `bson:"type" json:"type"` // "booking", "refund", "wallet_topup"
	Amount      float64            `bson:"amount" json:"amount"`
	BookingID   primitive.ObjectID `bson:"bookingId,omitempty" json:"bookingId,omitempty"` // может быть null для topup
	Status      string             `bson:"status" json:"status"`                           // "pending", "completed", "failed"
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
}
