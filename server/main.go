package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
	Name        string `json:"name" bson:"name"`
	Author      string `json:"author" bson:"author"`
	PublishYear int    `json:"publish_year" bson:"publish_year"`
}

type BooksApiResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Documents []Book `json:"docs,omitempty"`
}

type StatusApiResponse struct {
	Timestamp int64 `json:"timestamp"`
}

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

func getStatus(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().UTC().UnixMilli()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StatusApiResponse{Timestamp: timestamp})
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if book.Name == "" || book.Author == "" || book.PublishYear == 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := collection.InsertOne(r.Context(), book)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BooksApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	response := BooksApiResponse{Success: true, Message: "1 book inserted."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var results []Book

	cursor, err := collection.Find(r.Context(), bson.D{})
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BooksApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	defer cursor.Close(r.Context())

	for cursor.Next(r.Context()) {
		var result Book
		if err := cursor.Decode(&result); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(BooksApiResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		results = append(results, result)
	}

	response := BooksApiResponse{Success: true, Message: "Documents matching filter.", Documents: results}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
