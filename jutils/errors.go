package jutils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func ProcessError(block string, caughtError error) {
	fmt.Printf("***** Error in block: %s *****\n", block)
	fmt.Printf("***** Stamp: %s *****\n", FriendlyTimestamp())
	fmt.Println(caughtError)
	fmt.Println("***** End Error Report *****")
}

func ProcessHttpError(block string, caughtError error, eCode int, w http.ResponseWriter) {
	ProcessError(block, caughtError)
	http.Error(w, caughtError.Error(), eCode)
}

func ProcessCustomHttpError(block string, customError string, eCode int, w http.ResponseWriter) error {
	asError := errors.New(strings.ToLower(customError))
	ProcessHttpError(block, asError, eCode, w)
	return asError
}
