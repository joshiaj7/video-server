package testutil

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"video-server/module/entity"

	"github.com/stretchr/testify/assert"
)

var (
	ErrorUnexpected = entity.NewError("Unexpected error", http.StatusInternalServerError)
)

// revive:disable:cognitive-complexity,cyclomatic

func AssertErrorExAc(t *testing.T, expected interface{}, actual error) bool {
	t.Helper()

	switch value := expected.(type) {
	case error:
		return assert.ErrorIs(t, actual, value)
	case string:
		return assert.Equal(t, value, actual.Error())
	default:
		return assert.Nil(t, actual)
	}
}

type StructWithMapInterface interface {
	ToMap() map[string]interface{}
}

func AssertStructExAc(t *testing.T, expectedAbstract interface{}, actualStruct StructWithMapInterface) bool {
	t.Helper()

	if expectedAbstract == nil {
		return assert.Nil(t, actualStruct)
	}

	switch expected := expectedAbstract.(type) {
	case map[string]interface{}:
		if expected == nil {
			return assert.Nil(t, actualStruct)
		}
		actual := actualStruct.ToMap()
		messages := []string{}
		for key, expectedValue := range expected {
			actualValue, ok := actual[key]
			if ok {
				if actualValue != expectedValue {
					messages = append(messages, fmt.Sprintf("Key %s -> Expected(%s)-(%T) Actual(%s)-(%t)", key, expectedValue, expectedValue, actualValue, actualValue))
				}
			} else {
				messages = append(messages, fmt.Sprintf("Key %s -> Expected(%s)-(%T) Actual Not Found", key, expectedValue, expectedValue))
			}
		}

		if len(messages) > 0 {
			return assert.Fail(t, strings.Join(messages, "\n"))
		}
		return true
	default:
		if expected == nil {
			return assert.Nil(t, actualStruct)
		}
		return assert.Equal(t, expected, actualStruct)
	}
}
