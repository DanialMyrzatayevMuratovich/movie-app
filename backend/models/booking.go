package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Booking struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BookingNumber string             `bson:"bookingNumber" json:"bookingNumber"` // "BK-20260201-001234"
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	ShowtimeID    primitive.ObjectID `bson:"showtimeId" json:"showtimeId"`
	Seats         []BookingSeat      `bson:"seats" json:"seats"` // Embedded
	TotalAmount   float64            `bson:"totalAmount" json:"totalAmount"`
	Status        string             `bson:"status" json:"status"`   // "pending", "confirmed", "cancelled", "expired"
	Payment       Payment            `bson:"payment" json:"payment"` // Embedded
	QRCode        string             `bson:"qrCode" json:"qrCode"`
	ExpiresAt     time.Time          `bson:"expiresAt" json:"expiresAt"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type BookingSeat struct {
	Row    string  `bson:"row" json:"row"`
	Number int     `bson:"number" json:"number"`
	Price  float64 `bson:"price" json:"price"`
}

type Payment struct {
	Method        string    `bson:"method" json:"method"` // "card", "wallet", "cash"
	TransactionID string    `bson:"transactionId" json:"transactionId"`
	PaidAt        time.Time `bson:"paidAt,omitempty" json:"paidAt,omitempty"`
	Status        string    `bson:"status" json:"status"` // "pending", "completed", "failed"
}
