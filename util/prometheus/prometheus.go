package prometheus

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"

	"github.com/yukels/util/log"
	"github.com/yukels/util/system"
)

const (
	ComponentName   = "component"
	PodName         = "pod_ip"
	EnvironmentName = "environment"
)

var (
	ip                             = system.OutboundIP()
	env                            = os.Getenv("ENVIRONMENT")
	component                      = os.Getenv("COMPONENT")
	registry  *prometheus.Registry = nil
)

func Registry() *prometheus.Registry {
	if registry == nil {
		if component == "" {
			component = log.ProgramName
		}
		registry = prometheus.NewRegistry()

		MustRegister(registry,
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(collectors.WithGoCollections(collectors.GoRuntimeMemStatsCollection|collectors.GoRuntimeMetricsCollection)))
	}
	return registry
}

func MustRegisterDefault(co ...prometheus.Collector) {
	MustRegister(registry, co...)
}

func MustRegister(reg prometheus.Registerer, co ...prometheus.Collector) {
	labels := prometheus.Labels{ComponentName: component, PodName: ip}
	if env != "" {
		labels[EnvironmentName] = env
	}

	prometheus.WrapRegistererWith(labels, reg).MustRegister(co...)
}
