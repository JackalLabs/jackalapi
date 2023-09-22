package jutils

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Helpers
func shouldPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() { _ = recover() }()
	f()
	t.Errorf("did not panic")
}

type testCase struct {
	name        string
	varId       string
	fallBack    string
	shouldPanic bool
}

var (
	tt = []testCase{
		{
			name:        "env var exists",
			varId:       "MOCK_ENV_VAR",
			shouldPanic: false,
		},
		{
			name:        "env var doesn't exist",
			varId:       "I_DONT_EXIST",
			fallBack:    "fallback_var",
			shouldPanic: true,
		},
	}
)

// Env funcs
func TestLoadEnvVarOrFallback(t *testing.T) {
	r := require.New(t)
	mockEnvValue := "lupulella-2"
	t.Setenv("MOCK_ENV_VAR", mockEnvValue)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := LoadEnvVarOrFallback(tc.varId, tc.fallBack)
			fmt.Println(tc.varId)

			if len(tc.fallBack) > 0 {
				r.Equal(value, tc.fallBack)
			} else {
				r.Equal(value, mockEnvValue)
			}
		})
	}
}

func TestLoadEnvVarOrPanic(t *testing.T) {
	r := require.New(t)
	mockEnvValue := "lupulella-2"
	t.Setenv("MOCK_ENV_VAR", mockEnvValue)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.varId)
			if tc.shouldPanic {
				shouldPanic(t, func() { LoadEnvVarOrPanic(tc.varId) })
			} else {
				r.Equal(LoadEnvVarOrPanic(tc.varId), mockEnvValue)
			}
		})
	}
}

// Data funcs
func TestCloneBytes(t *testing.T) {
	r := require.New(t)

	byteBuffer := new(bytes.Buffer)
	byteArray := byteBuffer.Bytes()

	byteReader := bytes.NewReader(byteArray)
	byteArray2 := CloneBytes(byteReader)

	matches := reflect.DeepEqual(byteArray, byteArray2)
	r.Equal(matches, true)
}

// Date funcs
func TestFriendlyTimestamp(t *testing.T) {
	r := require.New(t)

	now := FriendlyTimestamp()
	const layout = "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, now)
	if err != nil {
		t.Error(err)
	}
	r.IsType(now, "")
	r.IsType(parsedTime, time.Time{})
}

// Error funcs
func TestProcessError(t *testing.T) {
	// TODO - add test
}

func TestProcessHttpError(t *testing.T) {
	// TODO - add test
}

func TestProcessCustomHttpError(t *testing.T) {
	// TODO - add test
}
