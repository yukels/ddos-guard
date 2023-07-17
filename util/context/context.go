package context

import (
	"context"
	"time"
)

type key int64

const (
	userKey key = iota
	tokenKey
	actualUserKey
)

// Context is an implementation of the stdlib Context with strongly-typed specific values
type Context struct {
	context.Context
}

// New returns an instance of context wrapping given context.
// If given context is already context - just returns it.
func New(ctx context.Context) Context {
	if ctx, ok := ctx.(Context); ok {
		return ctx
	}
	return Context{Context: ctx}
}

// Background returns a non-nil, empty Context. It is never canceled, has no
// values, and has no deadline. It is typically used by the main function,
// initialization, and tests, and as the top-level Context for incoming
// requests.
func Background() Context {
	return New(context.Background())
}

// WithCancel returns a copy of parent with a new Done channel. The returned
// context's Done channel is closed when the returned cancel function is called
// or when the parent context's Done channel is closed, whichever happens first.
//
// Canceling this context releases resouces associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func WithCancel(parent Context) (Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	return Context{Context: ctx}, cancel
}

// WithTimeout returns WithDeadline(parent, time.Now().Add(timeout)).
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete:
//
//	func slowOperationWithTimeout(ctx context.Context) (Result, error) {
//		ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
//		defer cancel()  // releases resources if slowOperation completes before timeout elapses
//		return slowOperation(ctx)
//	}
func WithTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(timeout))
	return Context{Context: ctx}, cancel
}

// WithDefaultTimeout creates a new context with timeout if it wasn't setup before
func WithDefaultTimeout(parent Context, timeout time.Duration) (Context, context.CancelFunc) {
	if _, ok := parent.Deadline(); ok {
		return WithCancel(parent)
	}
	return WithTimeout(parent, timeout)
}

// Value returns the value associated with this context for key
func (ctx Context) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}

// User returns the current user's ID if available
func (ctx Context) User() *string {
	if user, ok := ctx.Context.Value(userKey).(string); ok {
		return &user
	}
	return nil
}

// Token returns the Token object if available
func (ctx Context) Token() *string {
	if token, ok := ctx.Context.Value(tokenKey).(string); ok {
		return &token
	}
	return nil
}

// WithValue creates a child context with given key and value
func (ctx Context) WithValue(key, val interface{}) Context {
	return Context{Context: context.WithValue(ctx, key, val)}
}

// WithUser creates a child context with given User
func (ctx Context) WithUser(user string) Context {
	return Context{Context: context.WithValue(ctx, userKey, user)}
}

// WithActualUser creates a child context with given User
func (ctx Context) WithActualUser(user string) Context {
	return Context{Context: context.WithValue(ctx, actualUserKey, user)}
}

// WithToken creates a child context with given Token object
func (ctx Context) WithToken(token string) Context {
	return Context{Context: context.WithValue(ctx, tokenKey, token)}
}
