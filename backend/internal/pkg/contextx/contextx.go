package contextx

// Custom context package for backend
// This package provides utilities for managing context in the backend,
// including request-scoped data and utilities for handling context values.
// It is designed to be used with the Echo framework and integrates with
// the application's middleware and handlers.
import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Contextx interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
	ReqContext() *echo.Context
	WithTxn(ctx context.Context) Contextx
	GetTxn(*gorm.DB) *gorm.DB
}

type contextx struct {
	context.Context
	reqCtx *echo.Context
	txn    *gorm.DB
}

func NewContextx(ctx context.Context) Contextx {
	return &contextx{
		Context: ctx,
		txn:     nil,
	}
}

func NewWithRequestContext(req echo.Context) Contextx {
	return &contextx{
		Context: req.Request().Context(),
		reqCtx:  &req,
		txn:     nil,
	}
}

func Background() Contextx {
	return &contextx{
		Context: context.Background(),
		reqCtx:  nil,
		txn:     nil,
	}
}

// Deadline returns the deadline for the context
func (c *contextx) Deadline() (time.Time, bool) {
	return c.Context.Deadline()
}

// Done returns a channel that is closed when the context is done
func (c *contextx) Done() <-chan struct{} {
	return c.Context.Done()
}

// Err returns an error if the context is done
func (c *contextx) Err() error {
	return c.Context.Err()
}

func (c *contextx) ReqContext() *echo.Context {
	if c.reqCtx == nil {
		return nil
	}
	return c.reqCtx
}

func (c *contextx) WithTxn(ctx context.Context) Contextx {
	c.txn = c.txn.WithContext(ctx)
	return c

}

func (c *contextx) GetTxn(db *gorm.DB) *gorm.DB {
	if c.txn == nil {
		return db
	}
	return c.txn.WithContext(c.Context)
}
