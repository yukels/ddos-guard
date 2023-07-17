package log

import (
	stdlog "log"

	"github.com/yukels/util/context"
)

type logWriter struct {
	ctx context.Context
}

// Write io.Writer iterface implementation
func (l *logWriter) Write(p []byte) (n int, err error) {
	Log(l.ctx).Debug(string(p))
	return len(p), nil
}

// NewStdLogger wrapper on standard go's logger
func NewStdLogger() *stdlog.Logger {
	writer := &logWriter{ctx: context.Background()}
	return stdlog.New(writer, "", 0)
}
