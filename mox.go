package mox

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/paywithclay/mox/model"
	"github.com/paywithclay/mox/query"
	"github.com/paywithclay/mox/session"
)

// Client represents the main MongoDB ORM client
type Client struct {
	*mongo.Client
	database *mongo.Database
}

// Connect creates a new MongoDB client with the given URI and database name
func Connect(uri, dbName string, opts ...*options.ClientOptions) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	for _, opt := range opts {
		clientOpts = opt
	}

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client:   client,
		database: client.Database(dbName),
	}, nil
}

// Collection returns a query.Collection for the given model
func (c *Client) Collection(m model.Model) *query.Collection {
	return &query.Collection{
		Collection: c.database.Collection(m.Table()),
	}
}

// NewSession creates a new MongoDB session
func (c *Client) NewSession(ctx context.Context) (*session.Session, error) {
	sess, err := c.Client.StartSession()
	if err != nil {
		return nil, err
	}
	return session.New(ctx, sess), nil
}

// WithTransaction executes a transaction
func (c *Client) WithTransaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	sess, err := c.NewSession(ctx)
	if err != nil {
		return nil, err
	}
	defer sess.EndSession(ctx)
	return sess.WithTransaction(fn)
}
