package handlers

import (
	"net/http"
	"time"

	"server/internal/repos"

	"github.com/wellers/webapi-shared/types"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetStatus(c *gin.Context) {
	timestamp := time.Now().UTC().UnixMilli()

	c.JSON(http.StatusOK, types.StatusApiResponse{Timestamp: timestamp})
}

func InsertBook(c *gin.Context, repo repos.BookRepository) {
	var book types.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, types.BooksApiResponse{
			Success: false,
			Message: "Invalid request payload",
		})
		return
	}

	if book.Name == "" || book.Author == "" || book.PublishYear == 0 {
		c.JSON(http.StatusBadRequest, types.BooksApiResponse{
			Success: false,
			Message: "Missing required fields",
		})
		return
	}

	err := repo.InsertOne(c, book)
	if err != nil {
		c.JSON(http.StatusOK, types.BooksApiResponse{
			Success: false,
			Message: "Failed to insert book: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.BooksApiResponse{
		Success: true,
		Message: "1 book inserted.",
	})
}

func GetBooks(c *gin.Context, repo repos.BookRepository) {
	var results []types.Book

	cursor, err := repo.Find(c, types.BookFindFilter{})
	if err != nil {
		c.JSON(http.StatusOK, types.BooksApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var result types.Book
		if err := cursor.Decode(&result); err != nil {
			c.JSON(http.StatusOK, types.BooksApiResponse{
				Success: false,
				Message: err.Error(),
			})
			return
		}
		results = append(results, result)
	}

	c.JSON(http.StatusOK, types.BooksApiResponse{
		Success:   true,
		Message:   "Documents matching filter.",
		Documents: results,
	})
}

func DeleteBook(c *gin.Context, repo repos.BookRepository) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, types.BooksApiResponse{
			Success: false,
			Message: "Missing required fields",
		})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusOK, types.BooksApiResponse{
			Success: false,
			Message: "Failed to parse id: " + err.Error(),
		})
		return
	}

	err = repo.DeleteOne(c, types.BookDeleteFilter{Id: objectId})
	if err != nil {
		c.JSON(http.StatusOK, types.BooksApiResponse{
			Success: false,
			Message: "Failed to delete book: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.BooksApiResponse{
		Success: true,
		Message: "1 book deleted.",
	})
}
