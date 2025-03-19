package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"server/internal/handlers"
	"server/internal/middleware"
	"server/internal/repos"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	r := gin.Default()

	collection := connectMongo()
	defer disconnectMongo()

	repo := &repos.MongoBookRepository{Collection: collection}

	validToken := os.Getenv("VALID_TOKEN")
	if validToken == "" {
		log.Fatal("VALID_TOKEN is not set")
		return
	}

	registerRoutes(r, repo, validToken)

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

func registerRoutes(r *gin.Engine, repo *repos.MongoBookRepository, validToken string) {
	r.GET("/status", handlers.GetStatus)

	authMiddleware := middleware.AuthMiddleware(validToken)

	authorised := r.Group("/")
	authorised.Use(authMiddleware)
	{
		authorised.POST("/book", func(c *gin.Context) {
			handlers.InsertBook(c, repo)
		})
		authorised.GET("/books", func(c *gin.Context) {
			handlers.GetBooks(c, repo)
		})
		authorised.DELETE("/book/:id", func(c *gin.Context) {
			handlers.DeleteBook(c, repo)
		})
	}
}
