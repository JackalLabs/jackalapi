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

var tt = []testCase{
	{
		name:        "env var exists",
		varId:       "MOCK_ENV_VAR",
		fallBack:    "",
		shouldPanic: false,
	},
	{
		name:        "env var doesn't exist",
		varId:       "I_DONT_EXIST",
		fallBack:    "fallback_var",
		shouldPanic: true,
	},
}

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

func TestCloneByteSlice(t *testing.T) {
	r := require.New(t)

	byteBuffer := new(bytes.Buffer)
	byteArray := byteBuffer.Bytes()

	cloneArray1, cloneArray2, err := CloneByteSlice(byteArray)
	if err != nil {
		ProcessError("TestCloneByteSlice", err)
	}

	matches := reflect.DeepEqual(byteArray, cloneArray1)
	r.Equal(matches, true)
	matches = reflect.DeepEqual(byteArray, cloneArray2)
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

func TestUnixMsTimestamp(t *testing.T) {
	// TODO - add test
}

// Keyring funcs
func TestContextKeyring(t *testing.T) {
	r := require.New(t)
	keyring := MakeContextKeyring()

	baselineKey := keyring.UseKey("ReqUniquePath")
	keyring.AddKeyToRing("Testing")
	testingKey := keyring.UseKey("Testing")

	r.IsType(keyring, ContextKeyring{})
	r.IsType(baselineKey, "")
	r.IsType(testingKey, "")
	r.Equal(len(keyring), 2)
}

// Error funcs
func TestProcessError(t *testing.T) {
	// TODO - add test
}

func TestProcessCustomError(t *testing.T) {
	// TODO - add test
}

func TestProcessHttpError(t *testing.T) {
	// TODO - add test
}

func TestProcessCustomHttpError(t *testing.T) {
	// TODO - add test
}
