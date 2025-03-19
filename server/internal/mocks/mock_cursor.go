package mocks

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
)

type MockCursor struct {
	Data       []bson.M
	CurrentIdx int
}

func (m *MockCursor) Next(ctx context.Context) bool {
	if m.CurrentIdx < len(m.Data) {
		m.CurrentIdx++
		return true
	}
	return false
}

func (m *MockCursor) Decode(result interface{}) error {
	if m.CurrentIdx == 0 || m.CurrentIdx > len(m.Data) {
		return errors.New("no data to decode")
	}
	bsonBytes, _ := bson.Marshal(m.Data[m.CurrentIdx-1]) // Convert BSON to bytes
	return bson.Unmarshal(bsonBytes, result)             // Decode into result
}

func (m *MockCursor) Close(ctx context.Context) error {
	return nil
}
