package mocks

import (
	"context"

	"server/internal/repos"

	"github.com/wellers/webapi-shared/types"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) InsertOne(ctx context.Context, book types.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockBookRepository) Find(ctx context.Context, filter types.BookFindFilter) (cur repos.BookCursor, err error) {
	args := m.Called(ctx, filter)

	data := []bson.M{
		{"name": "Book 1", "author": "Author 1", "publish_year": 2022},
		{"name": "Book 2", "author": "Author 2", "publish_year": 2023},
	}

	return &MockCursor{Data: data, CurrentIdx: 0}, args.Error(1)
}

func (m *MockBookRepository) DeleteOne(ctx context.Context, filter types.BookDeleteFilter) error {
	args := m.Called(ctx, filter)
	return args.Error(0)
}
