package system

import (
	"os"
	"os/signal"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

// RegistryInterrupt registers interrupt handler
func RegistryInterrupt(ctx context.Context) context.Context {
	// capture interrupt signals from OS
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		s := <-sig
		log.Log(ctx).Debugf("signal[%s] received", s.String())
		cancel()
	}()
	return ctx
}

func WaitForSignal(ctx context.Context) {
	<-ctx.Done()
}
