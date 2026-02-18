package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func Connect() {
	// Tạo client, không cần context ở đây nữa
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("❌ Error creating client:", err)
		return
	}

	// Tạo context để kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Kiểm tra ping
	if err := client.Ping(ctx, nil); err != nil {
		fmt.Println("❌ Cannot connect to MongoDB:", err)
		return
	}

	fmt.Println("✅ Connected to MongoDB!")
	Client = client
	DB = client.Database("social")
}
