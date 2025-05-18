package session

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Session wraps a mongo.Session with additional context and transaction support
type Session struct {
	session mongo.Session
	ctx    context.Context
}

// New creates a new session with context
func New(ctx context.Context, mongoSession mongo.Session) *Session {
	return &Session{
		session: mongoSession,
		ctx:    ctx,
	}
}

// Context returns the session's context
func (s *Session) Context() context.Context {
	return s.ctx
}

// EndSession ends the session
func (s *Session) EndSession(ctx context.Context) {
	s.session.EndSession(ctx)
}

// WithTransaction executes a transaction function
func (s *Session) WithTransaction(fn func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	return s.session.WithTransaction(s.ctx, fn)
}
