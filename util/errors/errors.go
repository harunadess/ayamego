package util

import (
	"os"

	logger "github.com/jordanjohnston/ayamego/util/logger"
)

// FatalErrorHandler logs an error if present, and exits
// the program with an exit code of 1
func FatalErrorHandler(msg string, err error) {
	if err != nil {
		logger.Error(msg, err)
		os.Exit(1)
	}
}

// StandardErrorHandler logs an error if present
func StandardErrorHandler(msg string, err error) {
	if err != nil {
		logger.Error(msg, err)
	}
}
