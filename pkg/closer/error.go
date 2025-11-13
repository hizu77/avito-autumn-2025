package closer

import "github.com/pkg/errors"

var (
	ErrGroupNotFound = errors.New("group not found")
	ErrAlreadyClosed = errors.New("already closed")
)
