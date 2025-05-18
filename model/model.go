package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model defines the base interface for all models
type Model interface {
	CollectionName() string
}

// Base provides common fields and functionality for all models
type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// CollectionName returns the default collection name
func (b *Base) CollectionName() string {
	return "models"
}

// SetTimestamps updates the CreatedAt and UpdatedAt fields
func (b *Base) SetTimestamps() {
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
}
