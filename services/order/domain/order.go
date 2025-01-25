package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	Id      string    `bson:"_id"`
	Name    string    `bson:"name"`
	Version int       `bson:"version"`
	Created time.Time `bson:"created"`
}

func CreateOrder(ctx context.Context, id string, name string) (*Order, error) {
	if id == "" {
		id = uuid.NewString()
	}
	return &Order{
		Id:      id,
		Name:    name,
		Version: 1,
		Created: time.Now(),
	}, nil
}
