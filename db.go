package ubm

import (
	"errors"
)

const (
	collectionActionName = "ubm"
)

var (
	ErrActionNotFound = errors.New("action for id not found")
	ErrUserNotFound   = errors.New("id not found")
)

// DB is abstract database interface for different usage cases
type DB interface {
	AddAction(id interface{}, actionName string) error
	GetAction(id interface{}, actionName string) (Action, error)
	GetLastAction(id interface{}) (LastAction, error)
}
