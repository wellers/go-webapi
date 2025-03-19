package repos

import (
	"context"
	"errors"

	"server/internal/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookCursor interface {
	Next(ctx context.Context) bool
	Decode(result interface{}) error
	Close(ctx context.Context) error
}

type BookRepository interface {
	InsertOne(ctx context.Context, book types.Book) error
	Find(ctx context.Context, filter types.BookFindFilter) (cur BookCursor, err error)
	DeleteOne(ctx context.Context, filter types.BookDeleteFilter) error
}

type MongoBookRepository struct {
	Collection *mongo.Collection
}

func (m *MongoBookRepository) Find(ctx context.Context, filter types.BookFindFilter) (cur BookCursor, err error) {
	return m.Collection.Find(ctx, filter)
}

func (m *MongoBookRepository) InsertOne(ctx context.Context, book types.Book) error {
	book.Id = primitive.NewObjectID()
	result, err := m.Collection.InsertOne(ctx, book)

	if result.InsertedID == nil {
		return errors.New("failed to insert")
	}

	return err
}

func (m *MongoBookRepository) DeleteOne(ctx context.Context, filter types.BookDeleteFilter) error {
	result, err := m.Collection.DeleteOne(ctx, filter)

	if result.DeletedCount != 1 {
		return errors.New("failed to delete")
	}

	return err
}
