package guard

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/yukels/util/context"
	prometheusutil "github.com/yukels/util/prometheus"
)

const (
	userBlockedMetricName = "ddos_user_blocked"
	proxyGuardMetricName  = "ddos_guard_status"
)

type Collector struct {
	userBlocked *prometheus.CounterVec
	guardStatus *prometheus.GaugeVec
}

// NewCollector prometheus metrics collector
func NewCollector(ctx context.Context) *Collector {
	prometheusutil.Registry()

	c := &Collector{
		userBlocked: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: userBlockedMetricName,
				Help: "DDOS guard blocked users",
			},
			[]string{"code", "user"},
		),
		guardStatus: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: proxyGuardMetricName,
				Help: "DDOS guard status",
			},
			[]string{},
		),
	}
	prometheusutil.MustRegisterDefault(c.userBlocked, c.guardStatus)

	return c
}

func (c *Collector) IncUserBlock(ctx context.Context, code, user string) {
	c.userBlocked.WithLabelValues(code, user).Inc()
}

func (c *Collector) SetGuardStatus(ctx context.Context, status int) {
	c.guardStatus.WithLabelValues().Set(float64(status))
}
