package prometheus

// PrometheusConfig represents a Prometheus configuration
type PrometheusConfig struct {
	Url      string
	Username string
	Password string
	Timeout  int
}
