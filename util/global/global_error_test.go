package global

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/sirupsen/logrus"

	"github.com/yukels/util/context"
)

func TestHandleGlobalError(t *testing.T) {
	ctx := context.Background()

	errorLogged := false

	patchLog := gomonkey.ApplyFunc((*logrus.Entry).Fatal, func(entry *logrus.Entry, args ...interface{}) { errorLogged = true })
	defer patchLog.Reset()

	defer func() {
		if !errorLogged {
			t.Errorf("Error was not logged")
		}
	}()
	defer HandleGlobalError(ctx)

	panic("panic")
}
