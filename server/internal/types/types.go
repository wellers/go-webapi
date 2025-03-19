package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Author      string             `json:"author" bson:"author"`
	PublishYear int                `json:"publish_year" bson:"publish_year"`
}

type BooksApiResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Documents []Book `json:"docs,omitempty"`
}

type StatusApiResponse struct {
	Timestamp int64 `json:"timestamp"`
}
