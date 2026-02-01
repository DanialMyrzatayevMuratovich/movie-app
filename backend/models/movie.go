package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Movie struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title          string             `bson:"title" json:"title" validate:"required"`
	TitleKz        string             `bson:"titleKz" json:"titleKz"`
	TitleRu        string             `bson:"titleRu" json:"titleRu"`
	Description    string             `bson:"description" json:"description"`
	Director       string             `bson:"director" json:"director"`
	Cast           []string           `bson:"cast" json:"cast"`
	Genres         []string           `bson:"genres" json:"genres"`
	Duration       int                `bson:"duration" json:"duration"` // minutes
	ReleaseDate    time.Time          `bson:"releaseDate" json:"releaseDate"`
	Rating         string             `bson:"rating" json:"rating"` // "G", "PG", "PG-13", "R"
	IMDBRating     float64            `bson:"imdbRating" json:"imdbRating"`
	Language       []string           `bson:"language" json:"language"` // ["Kazakh", "Russian", "English"]
	Subtitles      []string           `bson:"subtitles" json:"subtitles"`
	PosterFileID   primitive.ObjectID `bson:"posterFileId,omitempty" json:"posterFileId,omitempty"` // GridFS
	PosterURL      string             `bson:"posterUrl,omitempty" json:"posterUrl,omitempty"`
	TrailerFileID  primitive.ObjectID `bson:"trailerFileId,omitempty" json:"trailerFileId,omitempty"` // GridFS
	IsActive       bool               `bson:"isActive" json:"isActive"`
	AgeRestriction int                `bson:"ageRestriction" json:"ageRestriction"`
	Reviews        []Review           `bson:"reviews" json:"reviews"` // Embedded
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
}

type Review struct {
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	Rating    float64            `bson:"rating" json:"rating"` // 1-10
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}
