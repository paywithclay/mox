package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Relationship struct {
	Type       string        // "belongsTo", "hasMany", "hasOne", "embeds"
	Collection string        
	Field      string        
	Pipeline   []bson.D      
	Options    *options.FindOptions
}

// BelongsTo defines a one-to-one relationship
func BelongsTo(collection string) *Relationship {
	return &Relationship{
		Type: "belongsTo", 
		Collection: collection,
	}
}

// HasMany defines a one-to-many relationship
func HasMany(collection string) *Relationship {
	return &Relationship{
		Type: "hasMany", 
		Collection: collection,
	}
}

// WithPipeline adds an aggregation pipeline to the relationship
func (r *Relationship) WithPipeline(pipeline []bson.D) *Relationship {
	r.Pipeline = pipeline
	return r
}

// WithOptions adds find options to the relationship
func (r *Relationship) WithOptions(opts *options.FindOptions) *Relationship {
	r.Options = opts
	return r
}

// Load executes the relationship query with pipeline support
func (r *Relationship) Load(ctx context.Context, db *mongo.Database, parentID interface{}) (*mongo.Cursor, error) {
	coll := db.Collection(r.Collection)
	
	if len(r.Pipeline) > 0 {
		// Add match stage for parent ID based on relationship type
		matchStage := bson.D{}
		switch r.Type {
		case "belongsTo":
			matchStage = bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: parentID}}}}
		case "hasMany":
			matchStage = bson.D{{Key: "$match", Value: bson.D{{Key: "parent_id", Value: parentID}}}}
		}
		
		fullPipeline := append([]bson.D{matchStage}, r.Pipeline...)
		return coll.Aggregate(ctx, fullPipeline)
	}
	
	switch r.Type {
	case "belongsTo":
		return coll.Find(ctx, bson.M{"_id": parentID}, r.Options)
	case "hasMany":
		return coll.Find(ctx, bson.M{"parent_id": parentID}, r.Options)
	case "embeds":
		return coll.Find(ctx, bson.M{"_id": parentID}, r.Options)
	default:
		return nil, nil
	}
}
