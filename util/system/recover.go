package system

import (
	"github.com/pkg/errors"
)

// RecoverError extracts error from recover
func RecoverError() error {
	if r := recover(); r != nil {
		// find out exactly what the error was and set err
		switch x := r.(type) {
		case string:
			return errors.New(x)
		case error:
			return x
		default:
			// Fallback err (per specs, error strings should be lowercase w/o punctuation
			return errors.New("unknown panic")
		}
	}
	return nil
}
