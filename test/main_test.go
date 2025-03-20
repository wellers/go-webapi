package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"test/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connectMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")

	var err error
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	// verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB...")

	return client
}

func disconnectMongo(client *mongo.Client) {
	// Wait for the server to finish and disconnect gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("MongoDB disconnect error:", err)
	}
	fmt.Println("MongoDB disconnected.")
}

type ServerTestSuite struct {
	suite.Suite
	client   *mongo.Client
	database *mongo.Database
}

var testId primitive.ObjectID

func (suite *ServerTestSuite) SetupSuite() {
	suite.client = connectMongo()
	suite.database = suite.client.Database("my_database")
}

func (suite *ServerTestSuite) SetupTest() {
	coll := suite.database.Collection("books")

	testId = primitive.NewObjectID()
	_, err := coll.InsertOne(context.Background(), types.Book{
		Id:          testId,
		Name:        "Alice's Adventures in Wonderland",
		Author:      "Lewis Carroll",
		PublishYear: 1865,
	})
	suite.NoError(err)
}

func (suite *ServerTestSuite) TearDownTest() {
	suite.database.Drop(context.Background())
}

func (suite *ServerTestSuite) TearDownSuite() {
	disconnectMongo(suite.client)
}

func (suite *ServerTestSuite) TestPostBook() {
	serverUrl := os.Getenv("SERVER_URI")
	validToken := os.Getenv("VALID_TOKEN")

	jsonData := `{
		"name": "War of the Worlds",
		"author": "H.G. Wells",
		"publish_year": 1898
	}`

	req, err := http.NewRequest("POST", serverUrl+"/api/v1/books", bytes.NewBuffer([]byte(jsonData)))
	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Add("Accept", "application/json")
	assert.NoError(suite.T(), err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var response types.BooksApiResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "1 book inserted.", response.Message)
}

func (suite *ServerTestSuite) TestGetBooks() {
	serverUrl := os.Getenv("SERVER_URI")
	validToken := os.Getenv("VALID_TOKEN")
	req, err := http.NewRequest("GET", serverUrl+"/api/v1/books", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Add("Accept", "application/json")
	assert.NoError(suite.T(), err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var response types.BooksApiResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "Documents matching filter.", response.Message)
	assert.Equal(suite.T(), []types.Book{{
		Id:          testId,
		Name:        "Alice's Adventures in Wonderland",
		Author:      "Lewis Carroll",
		PublishYear: 1865,
	}}, response.Documents)
}

func (suite *ServerTestSuite) TestDeleteBook() {
	serverUrl := os.Getenv("SERVER_URI")
	validToken := os.Getenv("VALID_TOKEN")

	req, err := http.NewRequest("DELETE", serverUrl+"/api/v1/books/"+testId.Hex(), bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+validToken)
	req.Header.Add("Accept", "application/json")
	assert.NoError(suite.T(), err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(suite.T(), err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(suite.T(), err)

	var response types.BooksApiResponse
	err = json.Unmarshal(body, &response)
	assert.NoError(suite.T(), err)

	assert.True(suite.T(), response.Success)
	assert.Equal(suite.T(), "1 book deleted.", response.Message)
}

func getEnvAsBool(name string, defaultValue bool) bool {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		fmt.Printf("Warning: Invalid boolean value for %s: %s\n", name, valStr)
		return defaultValue
	}
	return val
}

func TestServer(t *testing.T) {
	if getEnvAsBool("RUN_TESTS", false) {
		suite.Run(t, new(ServerTestSuite))
	}
}
