package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"server/internal/handlers"
	"server/internal/repos"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	collection := connectMongo()
	defer disconnectMongo()

	r := gin.Default()

	repo := &repos.MongoBookRepository{Collection: collection}
	registerRoutes(r, repo)

	r.Run(":80")
}

func connectMongo() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	// verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB...")

	return client.Database("my_database").Collection("books")
}

func disconnectMongo() {
	// Wait for the server to finish and disconnect gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("MongoDB disconnect error:", err)
	}
	fmt.Println("MongoDB disconnected.")
}

func registerRoutes(r *gin.Engine, repo *repos.MongoBookRepository) {
	r.GET("/status", handlers.GetStatus)
	r.POST("/book", func(c *gin.Context) {
		handlers.InsertBook(c, repo)
	})
	r.GET("/books", func(c *gin.Context) {
		handlers.GetBooks(c, repo)
	})
	r.DELETE("/book/:id", func(c *gin.Context) {
		handlers.DeleteBook(c, repo)
	})
}
