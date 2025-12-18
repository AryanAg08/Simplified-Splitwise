package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var GroupCollection *mongo.Collection
var ExpenseCollection *mongo.Collection
var BalanceCollection *mongo.Collection

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://localhost:27017/splitwise",
	))

	if err != nil {
		log.Fatal(" MongoDb Connection Failed!!", err)
	}

	// ping db
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("ping failed:", err)
	}

	log.Println("Mongo Connected!!")

	MongoClient = client
	// NotesCollection = client.Database("notesgo").Collection("notes")
	GroupCollection = client.Database("splitwise").Collection("groups")
	ExpenseCollection = client.Database("splitwise").Collection("expenses")
	BalanceCollection = client.Database("splitwise").Collection("balance")
}
