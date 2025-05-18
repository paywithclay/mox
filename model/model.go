package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model defines the base interface for all models
type Model interface {
	Table() string
}

// Document provides common fields and functionality for all models
type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

// Table returns the default collection name
func (d *Document) Table() string {
	return "documents"
}

// SetTimestamps updates the CreatedAt and UpdatedAt fields
func (d *Document) SetTimestamps() {
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now
}

// Archive soft deletes the document
func (d *Document) Archive() {
	now := time.Now()
	d.DeletedAt = &now
}

// Restore un-deletes the document
func (d *Document) Restore() {
	d.DeletedAt = nil
}
