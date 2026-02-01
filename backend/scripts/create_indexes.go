package scripts

import (
	"cinema-booking/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateIndexes() {
	ctx := context.Background()

	// 1. Users indexes
	usersCol := config.GetCollection("users")
	createIndex(ctx, usersCol, "email", true) // unique email

	// 2. Showtimes indexes (compound для поиска сеансов)
	showtimesCol := config.GetCollection("showtimes")
	createCompoundIndex(ctx, showtimesCol, []string{"movieId", "cinemaId", "startTime"})
	createIndex(ctx, showtimesCol, "startTime", false)

	// 3. Bookings indexes
	bookingsCol := config.GetCollection("bookings")
	createCompoundIndex(ctx, bookingsCol, []string{"userId", "createdAt"})
	createIndex(ctx, bookingsCol, "bookingNumber", true) // unique
	// TTL index для автоудаления expired bookings
	createTTLIndex(ctx, bookingsCol, "expiresAt", 0)

	// 4. Movies indexes (text search)
	moviesCol := config.GetCollection("movies")
	createTextIndex(ctx, moviesCol, []string{"title", "description"})
	createIndex(ctx, moviesCol, "isActive", false)

	// 5. Cinemas indexes (geospatial)
	cinemasCol := config.GetCollection("cinemas")
	create2DSphereIndex(ctx, cinemasCol, "location")

	// 6. Transactions indexes
	transactionsCol := config.GetCollection("transactions")
	createCompoundIndex(ctx, transactionsCol, []string{"userId", "createdAt"})

	log.Println("✅ All indexes created successfully")
}

// Helper functions
func createIndex(ctx context.Context, col *mongo.Collection, field string, unique bool) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetUnique(unique),
	}
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("⚠️ Warning: Index on %s.%s already exists or failed: %v", col.Name(), field, err)
	} else {
		log.Printf("✅ Created index on %s.%s", col.Name(), field)
	}
}

func createCompoundIndex(ctx context.Context, col *mongo.Collection, fields []string) {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: 1})
	}
	indexModel := mongo.IndexModel{Keys: keys}
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("⚠️ Warning: Compound index on %s failed: %v", col.Name(), err)
	} else {
		log.Printf("✅ Created compound index on %s: %v", col.Name(), fields)
	}
}

func createTextIndex(ctx context.Context, col *mongo.Collection, fields []string) {
	keys := bson.D{}
	for _, field := range fields {
		keys = append(keys, bson.E{Key: field, Value: "text"})
	}
	indexModel := mongo.IndexModel{Keys: keys}
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("⚠️ Warning: Text index on %s failed: %v", col.Name(), err)
	} else {
		log.Printf("✅ Created text index on %s", col.Name())
	}
}

func create2DSphereIndex(ctx context.Context, col *mongo.Collection, field string) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: field, Value: "2dsphere"}},
	}
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("⚠️ Warning: 2dsphere index on %s.%s failed: %v", col.Name(), field, err)
	} else {
		log.Printf("✅ Created 2dsphere index on %s.%s", col.Name(), field)
	}
}

func createTTLIndex(ctx context.Context, col *mongo.Collection, field string, expireAfterSeconds int32) {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: field, Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(expireAfterSeconds),
	}
	_, err := col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("⚠️ Warning: TTL index on %s.%s failed: %v", col.Name(), field, err)
	} else {
		log.Printf("✅ Created TTL index on %s.%s", col.Name(), field)
	}
}
