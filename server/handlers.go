package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func getStatus(c *gin.Context) {
	timestamp := time.Now().UTC().UnixMilli()

	c.JSON(http.StatusOK, StatusApiResponse{Timestamp: timestamp})
}

func insertBook(c *gin.Context, repo BookRepository) {
	var book Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, BooksApiResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	if book.Name == "" || book.Author == "" || book.PublishYear == 0 {
		c.JSON(http.StatusBadRequest, BooksApiResponse{
			Success: false,
			Message: "Missing required fields",
		})
		return
	}

	err := repo.InsertOne(c, book)
	if err != nil {
		c.JSON(http.StatusOK, BooksApiResponse{
			Success: false,
			Message: "Failed to insert book: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BooksApiResponse{
		Success: true,
		Message: "1 book inserted.",
	})
}

func getBooks(c *gin.Context, repo BookRepository) {
	var results []Book

	cursor, err := repo.Find(c, bson.D{})
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
