package prometheus

import (
	"github.com/pkg/errors"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

const (
	defaultTimeout = 30
)

var (
	client Provider
)

// Connect to the prometheus
func Connect(ctx context.Context, config *PrometheusConfig) error {
	if config.Url == "" {
		log.Log(ctx).Warn("Prometheus url is not defined.")
		return nil
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	var err error
	client, err = New(ctx, config)
	return err
}

func Instance(ctx context.Context) (Provider, error) {
	if client == nil {
		return nil, errors.Errorf("Prometheus client is not initialized")
	}
	return client, nil
}
