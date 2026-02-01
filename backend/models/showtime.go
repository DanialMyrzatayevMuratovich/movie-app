package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Showtime struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	MovieID        primitive.ObjectID `bson:"movieId" json:"movieId"`
	CinemaID       primitive.ObjectID `bson:"cinemaId" json:"cinemaId"`
	HallID         primitive.ObjectID `bson:"hallId" json:"hallId"`
	StartTime      time.Time          `bson:"startTime" json:"startTime"`
	EndTime        time.Time          `bson:"endTime" json:"endTime"`
	BasePrice      float64            `bson:"basePrice" json:"basePrice"`
	Format         string             `bson:"format" json:"format"`     // "2D", "3D", "IMAX"
	Language       string             `bson:"language" json:"language"` // "Russian", "English", "Kazakh"
	Subtitles      string             `bson:"subtitles" json:"subtitles"`
	AvailableSeats int                `bson:"availableSeats" json:"availableSeats"`
	BookedSeats    []BookedSeat       `bson:"bookedSeats" json:"bookedSeats"`
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
}

type BookedSeat struct {
	Row    string `bson:"row" json:"row"`
	Number int    `bson:"number" json:"number"`
	Status string `bson:"status" json:"status"` // "available", "booked", "reserved"
}
