package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hall struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CinemaID   primitive.ObjectID `bson:"cinemaId" json:"cinemaId"`
	Name       string             `bson:"name" json:"name"`
	HallNumber int                `bson:"hallNumber" json:"hallNumber"`
	Capacity   int                `bson:"capacity" json:"capacity"`
	Type       string             `bson:"type" json:"type"` // "Standard", "VIP", "IMAX", "4DX"
	Seats      []Seat             `bson:"seats" json:"seats"`
}

type Seat struct {
	Row    string  `bson:"row" json:"row"`       // "A", "B", "C"
	Number int     `bson:"number" json:"number"` // 1, 2, 3...
	Type   string  `bson:"type" json:"type"`     // "regular", "vip", "couple"
	Price  float64 `bson:"price" json:"price"`   // базовая цена
}
