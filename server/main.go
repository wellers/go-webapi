package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var client *mongo.Client

func main() {
	collection = connectMongo()
	defer disconnectMongo()

	r := mux.NewRouter()
	registerRoutes(r)
	listenAndServe(r)
}

func connectMongo() *mongo.Collection {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://192.168.50.100:27017"))
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

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/status", getStatus).Methods("GET")
	r.HandleFunc("/book", insertBook).Methods("POST")
	r.HandleFunc("/books", getBooks).Methods("GET")
}

func listenAndServe(r *mux.Router) {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}

	fmt.Println("Listening on port 80...")

	err = http.Serve(listener, r)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
