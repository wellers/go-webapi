package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/handlers"
	"server/internal/mocks"

	"github.com/wellers/webapi-shared/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStatus_Success(t *testing.T) {
	router := gin.Default()
	router.GET("/status", handlers.GetStatus)

	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedTime := time.Now().UTC().UnixMilli()
	var response map[string]int64
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	actualTime, exists := response["timestamp"]
	assert.True(t, exists)
	assert.InDelta(t, expectedTime, actualTime, 100)
}

func TestInsertBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockBookRepository)

	mockRepo.On("InsertOne", mock.Anything, mock.Anything).Return(nil, nil)

	router.POST("/api/v1/books", func(c *gin.Context) {
		handlers.InsertBook(c, mockRepo)
	})

	body := `{"name": "The Go Programming Language", "author": "Alan Donovan", "publish_year": 2015}`
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "applications/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.BooksApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "1 book inserted.", response.Message)

	mockRepo.AssertExpectations(t)
}

func TestFind_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockBookRepository)

	mockRepo.On("Find", mock.Anything, mock.Anything).Return(nil, nil)

	router.GET("/api/v1/books", func(c *gin.Context) {
		handlers.GetBooks(c, mockRepo)
	})

	req, _ := http.NewRequest("GET", "/api/v1/books", nil)
	req.Header.Set("Content-Type", "applications/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.BooksApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Documents matching filter.", response.Message)
	assert.Equal(t, []types.Book{
		{Name: "Book 1", Author: "Author 1", PublishYear: 2022},
		{Name: "Book 2", Author: "Author 2", PublishYear: 2023},
	}, response.Documents)

	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	mockRepo := new(mocks.MockBookRepository)

	mockRepo.On("DeleteOne", mock.Anything, mock.Anything).Return(nil, nil)

	router.DELETE("/api/v1/books/:id", func(c *gin.Context) {
		handlers.DeleteBook(c, mockRepo)
	})

	req, _ := http.NewRequest("DELETE", "/api/v1/books/67b62e32b475b5eccc413498", nil)
	req.Header.Set("Content-Type", "applications/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response types.BooksApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "1 book deleted.", response.Message)

	mockRepo.AssertExpectations(t)
}
