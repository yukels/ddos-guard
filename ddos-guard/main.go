package main

import (
	"github.com/yukels/ddos-guard/config"
	"github.com/yukels/ddos-guard/guard"
	"github.com/yukels/util/context"
	"github.com/yukels/util/global"
	"github.com/yukels/util/log"
	"github.com/yukels/util/prometheus"
	"github.com/yukels/util/system"
)

func main() {
	ctx := context.Background()
	defer log.HandleExit(ctx)

	ctx = system.RegistryInterrupt(ctx)
	ctx = log.RegisterFatal(ctx)

	log.Init("ddos-guard")

	defer global.HandleGlobalError(ctx)

	version := "0.0.0"
	log.Log(ctx).Infof("ddos-guard version [%s]", version)

	if err := prometheus.Connect(ctx, &config.Configs.DdosGuardConfig.Prometheus); err != nil {
		log.Log(ctx).WithError(err).Fatal("Error connecting to prometheus")
	}

	guard, err := guard.NewDdosProxy(ctx)
	if err != nil {
		log.Log(ctx).WithError(err).Fatalf("Can't create proxy")
	}

	if err := guard.Run(ctx); err != nil {
		log.Log(ctx).WithError(err).Fatalf("Can't start proxy")
	}

	log.Log(ctx).Info("Finished")
}
