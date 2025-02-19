package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getStatus(c *gin.Context) {
	timestamp := time.Now().UTC().UnixMilli()

	c.JSON(http.StatusOK, StatusApiResponse{Timestamp: timestamp})
}

func insertBook(c *gin.Context) {
	var book Book
	if err := json.NewDecoder(c.Request.Body).Decode(&book); err != nil {
		http.Error(c.Writer, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if book.Name == "" || book.Author == "" || book.PublishYear == 0 {
		http.Error(c.Writer, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := collection.InsertOne(c, book)
	if err != nil {
		c.JSON(http.StatusOK, BooksApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BooksApiResponse{Success: true, Message: "1 book inserted."})
}

func getBooks(c *gin.Context) {
	var results []Book

	cursor, err := collection.Find(c, bson.D{})
	if err != nil {
		c.JSON(http.StatusOK, BooksApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var result Book
		if err := cursor.Decode(&result); err != nil {
			c.JSON(http.StatusOK, BooksApiResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, BooksApiResponse{Success: true, Message: "Documents matching filter.", Documents: results})
}
