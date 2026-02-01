package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cinema struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name         string               `bson:"name" json:"name" validate:"required"`
	City         string               `bson:"city" json:"city" validate:"required"`
	Address      string               `bson:"address" json:"address" validate:"required"`
	Location     Location             `bson:"location" json:"location"`
	Facilities   []string             `bson:"facilities" json:"facilities"` // ["3D", "IMAX", "VIP", "Parking"]
	HallIDs      []primitive.ObjectID `bson:"hallIds" json:"hallIds"`       // Referenced
	Rating       float64              `bson:"rating" json:"rating"`
	TotalReviews int                  `bson:"totalReviews" json:"totalReviews"`
	Images       []string             `bson:"images" json:"images"` // пути к файлам
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"`
}

type Location struct {
	Type        string    `bson:"type" json:"type"`               // "Point"
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}
