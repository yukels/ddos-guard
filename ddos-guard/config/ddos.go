package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
	"github.com/yukels/util/prometheus"
)

var (
	defaultHostOut    = "localhost"
	defaultPortIn     = 8081
	defaultPortOut    = 8080
	defaultRetryAfter = 60
)

type UserServiceConfig struct {
	RefreshPeriod    time.Duration
	WhiteListUsers   []string
	BlockedListUsers []string
	S3Bucket         string
	S3Path           string
}

type ProxyConfig struct {
	PortIn      int
	PortOut     int
	HostOut     string
	RetryAfter  int
	HealthPath  string
	MetricsPath string
}

type PrometheusQuery struct {
	Query      string
	UpperBound float64
}

type CloudWatchQuery struct {
	Namespace      string
	Metric         string
	DimensionName  string
	DimensionValue string
	UpperBound     float64
}

type MonitoringConfig struct {
	MetricsPeriodSeconds int
	PrometheusQueries    map[string]PrometheusQuery
	CloudWatchQueries    map[string]CloudWatchQuery
}

type GuardConfig struct {
	BucketDuration  time.Duration
	BucketsHistory  int
	TopUserCount    int
	FilterRatioStep int64
}

type DdosGuardConfig struct {
	Proxy       ProxyConfig
	UserService UserServiceConfig
	Monitoring  MonitoringConfig
	Guard       GuardConfig
	Prometheus  prometheus.PrometheusConfig
}

func (c *ProxyConfig) ReadFromEnv(ctx context.Context) error {
	c.defaults(ctx)

	if err := parseEnvironmentInt(ctx, "PORT_IN", c.PortIn, &c.PortIn); err != nil {
		return err
	}

	if err := parseEnvironmentInt(ctx, "PORT_OUT", c.PortOut, &c.PortOut); err != nil {
		return err
	}

	if err := parseEnvironmentInt(ctx, "RETRY_AFTER", c.RetryAfter, &c.RetryAfter); err != nil {
		return err
	}

	hostOut := os.Getenv("HOST_OUT")
	if hostOut == "" {
		c.HostOut = hostOut
		log.Log(ctx).Infof("'HOST_OUT' is not defined. Continue with default [%s]", defaultHostOut)
	}

	return nil
}

func (c *ProxyConfig) defaults(ctx context.Context) {
	if c.PortIn == 0 {
		c.PortIn = defaultPortIn
	}
	if c.PortOut == 0 {
		c.PortOut = defaultPortOut
	}
	if c.RetryAfter == 0 {
		c.RetryAfter = defaultRetryAfter
	}
	if c.HostOut == "" {
		c.HostOut = defaultHostOut
	}
}

func parseEnvironmentInt(ctx context.Context, envName string, defaultValue int, value *int) error {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		log.Log(ctx).Infof("'%s' is not defined. Continue with default [%d]", envName, defaultValue)
		*value = defaultValue
		return nil
	}

	num, err := strconv.Atoi(fromEnv)
	if err != nil {
		return errors.Wrapf(err, "Can't get '%s' from environment, value [%s] is not number", envName, fromEnv)
	}
	*value = num
	return nil
}
