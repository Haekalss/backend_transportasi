package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	mongoString := os.Getenv("MONGO_URI") // dibaca setelah godotenv.Load()
	dbName := os.Getenv("MONGO_DB")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoString))
	if err != nil {
		fmt.Println("NewClient Error:", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		fmt.Println("Connect Error:", err)
		return
	}

	DB = client.Database(dbName)
	fmt.Println("MongoDB connected to:", dbName)
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
