package testutil

import (
	"database/sql/driver"
	"time"
)

// Mock Struct allowing sqlmock argument to receive any time.
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface.
func (AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// Mock Struct allowing sqlmock argument to receive any string.
type AnyString struct{}

// Match satisfies sqlmock.Argument interface.
func (AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
