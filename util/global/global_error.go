package global

import (
	"github.com/pkg/errors"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

// HandleGlobalError recovers from an error, making sure to close connections and other static objects
func HandleGlobalError(ctx context.Context) {
	if r := recover(); r != nil {
		err := ExtractError(r)
		log.Log(ctx).WithError(err).Fatal("Error occurred in the global call")
	}
}

// ExtractError cast recover to known error type
func ExtractError(r interface{}) error {
	var err error
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("unknown panic")
	}
	return err
}
