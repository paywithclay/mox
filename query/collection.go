package query

import (
	"context"
	"reflect"

	"github.com/paywithclay/mox/hooks"
	"github.com/paywithclay/mox/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// hasZeroID checks if the model has an empty ObjectID
func hasZeroID(model any) bool {
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")

	if idField.IsValid() && idField.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		return idField.Interface().(primitive.ObjectID).IsZero()
	}
	return true
}

// getID retrieves the _id field from the model
func getID(model any) any {
	v := reflect.ValueOf(model).Elem()
	idField := v.FieldByName("ID")

	if idField.IsValid() && idField.Type() == reflect.TypeOf(primitive.ObjectID{}) {
		return idField.Interface().(primitive.ObjectID)
	}
	return nil
}

// DB provides high-level database operations
type DB struct {
	db *mongo.Database
}

// Collection returns a MongoDB collection for the given model
func (db *DB) Collection(m model.Model) *mongo.Collection {
	return db.db.Collection(m.CollectionName())
}

// Save inserts or updates a model
func (db *DB) Save(ctx context.Context, model model.Model) error {
	coll := db.Collection(model)

	// Call BeforeSave hook if available
	if err := hooks.CallBeforeSave(model); err != nil {
		return err
	}

	// If ID is empty, insert new doc
	if hasZeroID(model) {
		_, err := coll.InsertOne(ctx, model)
		return err
	}

	// Else update existing
	update := bson.M{"$set": model}
	_, err := coll.UpdateByID(ctx, getID(model), update)
	return err
}

// FindByID loads a model by ID
func (db *DB) FindByID(ctx context.Context, model model.Model, id interface{}) error {
	coll := db.Collection(model)
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(model)

	if err == nil {
		_ = hooks.CallAfterFind(model)
	}

	return err
}

// Delete removes a model by ID
func (db *DB) Delete(ctx context.Context, model model.Model, id interface{}) error {
	coll := db.Collection(model)
	res, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
