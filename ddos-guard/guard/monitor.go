package guard

import (
	"os"
	"time"

	"github.com/yukels/ddos-guard/config"
	awsclient "github.com/yukels/util/aws-client"
	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
	"github.com/yukels/util/prometheus"
)

const (
	guardStatusOK = 0
)

type Monitor struct {
	config           *config.MonitoringConfig
	prometheus       prometheus.Provider
	collector        *Collector
	statuses         map[string]bool
	inHighUsage      bool
	queryToEnum      map[string]int
	cloudWatchClient *awsclient.CloudWatchClient
}

func NewMonitor(ctx context.Context, collector *Collector, config *config.MonitoringConfig) (*Monitor, error) {
	prometheus, err := prometheus.Instance(ctx)
	if err != nil {
		log.Log(ctx).WithError(err).Warn("Can't create Prometheus client")
	}

	cloudWatchClient, err := awsclient.NewCloudWatchClient(ctx)
	if err != nil {
		log.Log(ctx).WithError(err).Warnf("Can't create cloudwatch AWS client")
	}

	m := &Monitor{
		config:           config,
		collector:        collector,
		prometheus:       prometheus,
		statuses:         map[string]bool{},
		queryToEnum:      map[string]int{},
		inHighUsage:      false,
		cloudWatchClient: cloudWatchClient,
	}

	for name, query := range config.PrometheusQueries {
		query.Query = os.ExpandEnv(query.Query)
		config.PrometheusQueries[name] = query
		m.statuses[name] = true
	}
	for name, query := range config.CloudWatchQueries {
		query.DimensionValue = os.ExpandEnv(query.DimensionValue)
		config.CloudWatchQueries[name] = query
		m.statuses[name] = true
	}

	return m, nil
}

func (m *Monitor) Run(ctx context.Context) error {
	m.collectMetrics(ctx)
	go m.metricsLoop(ctx)
	return nil
}

func (m *Monitor) InHighUsage(ctx context.Context) bool {
	return m.inHighUsage
}

func (m *Monitor) calcHighUsage(ctx context.Context) bool {
	for name, status := range m.statuses {
		if !status {
			log.Log(ctx).Warnf("InHighUsage due to [%s]", name)
			m.collector.SetGuardStatus(ctx, m.queryToEnum[name])
			return true
		}
	}
	m.collector.SetGuardStatus(ctx, guardStatusOK)
	return false
}

func (m *Monitor) metricsLoop(ctx context.Context) {
	log.Log(ctx).Info("Monitor thread is running...")
	waitPeriod := time.Duration(m.config.MetricsPeriodSeconds) * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitPeriod):
			m.collectMetrics(ctx)
		}
	}
}

func (m *Monitor) collectMetrics(ctx context.Context) {
	m.collectPrometheusMetrics(ctx)
	m.collectCloudWatchMetrics(ctx)
	m.inHighUsage = m.calcHighUsage(ctx)
}

func (m *Monitor) collectPrometheusMetrics(ctx context.Context) {
	if m.prometheus == nil {
		return
	}

	for name, query := range m.config.PrometheusQueries {
		result, err := m.prometheus.Query(ctx, query.Query)

		// set the status to 'true' on error - don't stop the guard due to the collector problem
		m.statuses[name] = true
		if err != nil {
			log.Log(ctx).WithError(err).Errorf("Prometheus error on [%s]", query.Query)
			continue
		}

		if len(result) == 0 {
			log.Log(ctx).Warnf("[%s] Got empty result for query", name)
			continue
		}

		value := float64(result[0].Value)
		log.Log(ctx).Debugf("[%s] Got Prometheus response value [%f]", name, value)
		if query.UpperBound < value {
			log.Log(ctx).Warnf("[%s] metric overloaded [%f < %f]", name, query.UpperBound, value)
			m.statuses[name] = false
		}
	}
}

func (m *Monitor) collectCloudWatchMetrics(ctx context.Context) {
	if m.cloudWatchClient == nil {
		return
	}

	for name, query := range m.config.CloudWatchQueries {
		value, err := m.cloudWatchClient.GetMetricLast(ctx, query.Namespace, query.Metric, query.DimensionName, query.DimensionValue)
		m.statuses[name] = true
		if err != nil {
			log.Log(ctx).WithError(err).Errorf("CloudWatch error on [%v]", query)
			continue
		}

		log.Log(ctx).Debugf("[%s] Got CloudWatch response value [%f]", name, value)
		if query.UpperBound < value {
			log.Log(ctx).Warnf("[%s] metric overloaded [%f < %f]", name, query.UpperBound, value)
			m.statuses[name] = false
		}
	}
}
