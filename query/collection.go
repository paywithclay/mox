package query

import (
	"context"
	"errors"
	"reflect"

	"github.com/paywithclay/mox/hooks"
	"github.com/paywithclay/mox/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Add type alias for Validatable
type Validatable = model.Validatable

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
	return db.db.Collection(m.Table())
}

// Save inserts or updates a model
func (db *DB) Save(ctx context.Context, model model.Model) error {
	// Call validation first
	if validatable, ok := model.(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

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

// Find retrieves a document by ID
func (db *DB) Find(ctx context.Context, model model.Model, id interface{}) error {
	return db.FindByID(ctx, model, id)
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

// Publish saves a document to the database
func (db *DB) Publish(ctx context.Context, doc model.Model) error {
	// Call validation first
	if validatable, ok := doc.(Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	if err := hooks.CallBeforeSave(doc); err != nil {
		return err
	}

	coll := db.Collection(doc)
	if hasZeroID(doc) {
		_, err := coll.InsertOne(ctx, doc)
		return err
	}

	update := bson.M{"$set": doc}
	_, err := coll.UpdateByID(ctx, getID(doc), update)
	return err
}

// Archive soft deletes a document
func (db *DB) Archive(ctx context.Context, doc model.Model) error {
	if d, ok := doc.(interface{ Archive() }); ok {
		d.Archive()
		return db.Publish(ctx, doc)
	}
	return errors.New("document does not support archiving")
}

// Restore un-deletes a document
func (db *DB) Restore(ctx context.Context, doc model.Model) error {
	if d, ok := doc.(interface{ Restore() }); ok {
		d.Restore()
		return db.Publish(ctx, doc)
	}
	return errors.New("document does not support restoring")
}

// Where adds a filter condition
func (db *DB) Where(field string, op string, value interface{}) *QueryBuilder {
	return &QueryBuilder{
		db:      db,
		filters: []bson.M{{field: bson.M{"$" + op: value}}},
	}
}

// SortBy specifies the sort order
func (db *DB) SortBy(field string, direction string) *QueryBuilder {
	return &QueryBuilder{
		db:  db,
		sort: bson.M{field: direction},
	}
}

// Limit sets the maximum number of results
func (db *DB) Limit(count int) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		limit: count,
	}
}

// QueryBuilder builds complex queries
type QueryBuilder struct {
	db      *DB
	filters []bson.M
	sort    bson.M
	limit   int
	skip    int
}

// Get executes the query and returns results
func (qb *QueryBuilder) Get(ctx context.Context, model model.Model) ([]interface{}, error) {
	coll := qb.db.Collection(model)
	filter := combineFilters(qb.filters)
	opts := options.Find().
		SetSort(qb.sort).
		SetLimit(int64(qb.limit)).
		SetSkip(int64(qb.skip))
	
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	
	var results []interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func combineFilters(filters []bson.M) bson.M {
	combined := bson.M{}
	for _, filter := range filters {
		for key, value := range filter {
			combined[key] = value
		}
	}
	return combined
}
