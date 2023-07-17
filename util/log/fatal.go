package log

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yukels/util/context"
)

var isFatal = false

// RegisterFatal replaces exit func by context cancelling
func RegisterFatal(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	handler := func() {
		cancel()
		os.Exit(1)
	}
	log.RegisterExitHandler(handler)
	return ctx
}

// HandleExit handles recover which logs the panic message
func HandleExit(ctx context.Context) {
	if err := recover(); err != nil {
		Log(ctx).WithField("err", err).Panic("Error occurred in the global call")
		isFatal = true
	}

	if isFatal {
		os.Exit(1)
	}
}
