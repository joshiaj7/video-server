package testutil

import (
	"errors"
	"time"
)

var (
	DB, DBMock = NewDatabase()

	CreatedAt = time.Now().Add(-1 * time.Hour)
	UpdatedAt = time.Now()

	ErrDB      = errors.New("DB Error")
	ErrStorage = errors.New("storage error")
)
