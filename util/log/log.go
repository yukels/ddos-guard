package log

import (
	"os"
	"runtime"

	"github.com/yukels/util/context"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	logRetentionDays int = 3
	defaultLogLevel      = log.InfoLevel
	IDField              = "id"
	CtxField             = "ctx"
)

var (
	ProgramName = ""
	cliMode     = false
)

// Init initializes log using specified programName
func Init(programName string) {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&formatter{})
	log.SetLevel(LogLevel())

	logger := log.StandardLogger()
	logger.ExitFunc = func(code int) {}

	ProgramName = programName

	// display stack error
	callerLevels := log.AllLevels
	stackLevels := []log.Level{log.PanicLevel, log.FatalLevel}
	log.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))
}

// SetDebugLevel used for verbodse mode
func SetDebugLevel() {
	log.SetLevel(log.DebugLevel)
}

// IsDebugLevel
func IsDebugLevel() bool {
	return logrus.GetLevel() >= logrus.DebugLevel
}

// RegisterExitHandler adds an exit handler
func RegisterExitHandler(handler func()) {
	log.RegisterExitHandler(handler)
}

// Log appends line, file and function context to the logger
func Log(ctx context.Context) *log.Entry {
	entry := log.NewEntry(log.StandardLogger())
	if !cliMode {
		if pc, file, line, ok := runtime.Caller(1); ok {
			fName := runtime.FuncForPC(pc).Name()
			entry = entry.WithField("file", file).WithField("line", line).WithField("func", fName)
		}
	}

	entry = entry.WithField(CtxField, ctx)

	return entry
}

func LogLevel() log.Level {
	level := os.Getenv("LOG_LEVEL")
	if len(level) == 0 {
		return defaultLogLevel
	}
	loglevel, err := log.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	return loglevel
}
