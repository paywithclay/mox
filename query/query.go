package query

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection provides enhanced query operations for MongoDB collections
type Collection struct {
	*mongo.Collection
}

// FindOneByID finds a document by its ID
func (c *Collection) FindOneByID(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.FindOne(ctx, bson.M{"_id": id}, opts...)
}

// UpdateOneByID updates a document by its ID
func (c *Collection) UpdateOneByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.UpdateOne(ctx, bson.M{"_id": id}, update, opts...)
}

// DeleteOneByID deletes a document by its ID
func (c *Collection) DeleteOneByID(ctx context.Context, id interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.DeleteOne(ctx, bson.M{"_id": id}, opts...)
}

// Exists checks if a document with the given ID exists
func (c *Collection) Exists(ctx context.Context, id interface{}, opts ...*options.FindOneOptions) (bool, error) {
	err := c.FindOneByID(ctx, id, opts...).Err()
	if err == nil {
		return true, nil
	}
	if err == mongo.ErrNoDocuments {
		return false, nil
	}
	return false, err
}
