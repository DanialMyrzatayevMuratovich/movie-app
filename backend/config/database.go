package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var MongoClient *mongo.Client

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Подключение к MongoDB
	clientOptions := options.Client().ApplyURI(AppConfig.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	// Проверка подключения
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	MongoClient = client
	DB = client.Database(AppConfig.MongoDatabase)

	log.Printf("✅ Connected to MongoDB database: %s", AppConfig.MongoDatabase)
}

func DisconnectDB() {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatal("❌ Error disconnecting from MongoDB:", err)
		}
		log.Println("✅ Disconnected from MongoDB")
	}
}

// Helper функции для получения коллекций
func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}
