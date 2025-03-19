package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCursor interface {
	Next(ctx context.Context) bool
	Decode(result interface{}) error
	Close(ctx context.Context) error
}

type BookRepository interface {
	InsertOne(ctx context.Context, book Book) error
	Find(ctx context.Context, filter bson.M) (cur BookCursor, err error)
	DeleteOne(ctx context.Context, filter bson.M) error
}

type MongoBookRepository struct {
	Collection *mongo.Collection
}

func (m *MongoBookRepository) Find(ctx context.Context, filter bson.M) (cur BookCursor, err error) {
	return m.Collection.Find(ctx, filter)
}

func (m *MongoBookRepository) InsertOne(ctx context.Context, book Book) error {
	_, err := m.Collection.InsertOne(ctx, book)
	return err
}

func (m *MongoBookRepository) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := m.Collection.DeleteOne(ctx, filter)
	return err
}
