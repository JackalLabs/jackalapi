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

// Env funcs
func TestLoadEnvVarOrFallback(t *testing.T) {
	t.Setenv("MOCK_ENV_VAR", "lupulella-2")

	tt := []struct {
		name     string
		varId    string
		fallBack string
	}{
		{
			name:  "env var exists",
			varId: "MOCK_ENV_VAR",
		},
		{
			name:     "env var doesn't exist",
			varId:    "I_DONT_EXIST",
			fallBack: "fallback_var",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value := LoadEnvVarOrFallback(tc.varId, tc.fallBack)
			fmt.Println(tc.varId)

			if len(tc.fallBack) > 0 {
				require.Equal(t, value, "fallback_var")
			} else {
				require.Equal(t, value, "lupulella-2")
			}
		})
	}
}

func TestLoadEnvVarOrPanic(t *testing.T) {
	t.Setenv("MOCK_ENV_VAR", "lupulella-2")

	tt := []struct {
		name  string
		varId string
	}{
		{
			name:  "env var exists",
			varId: "MOCK_ENV_VAR",
		},
		{
			name:  "env var doesn't exist",
			varId: "I_DONT_EXIST",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println(tc.varId)
			shouldPanic(t, func() { LoadEnvVarOrPanic(tc.varId) })
		})
	}
}

// Data funcs
func TestCloneBytes(t *testing.T) {
	byteBuffer := new(bytes.Buffer)
	byteArray := byteBuffer.Bytes()

	byteReader := bytes.NewReader(byteArray)
	byteArray2 := CloneBytes(byteReader)

	matches := reflect.DeepEqual(byteArray, byteArray2)
	require.Equal(t, matches, true)
}

// Date funcs
func TestFriendlyTimestamp(t *testing.T) {
	now := FriendlyTimestamp()
	parsedTime, err := time.Parse("2000-12-24 15:04:05", now)
	if err != nil {
		t.Error(err)
	}
	require.IsType(t, now, reflect.TypeOf(""))
	require.IsType(t, parsedTime, time.Time{})
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
