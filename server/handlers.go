package main

import (
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

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
