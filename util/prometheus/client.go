package prometheus

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	prometheusconfig "github.com/prometheus/common/config"
	"github.com/prometheus/common/model"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

type Client struct {
	client v1.API
	config *PrometheusConfig
}

func New(ctx context.Context, conf *PrometheusConfig) (Provider, error) {
	var authRoundTripper http.RoundTripper
	if conf.Username != "" {
		authRoundTripper = prometheusconfig.NewBasicAuthRoundTripper(conf.Username, prometheusconfig.Secret(conf.Password), "", api.DefaultRoundTripper)
	}

	client, err := api.NewClient(api.Config{
		Address:      conf.Url,
		RoundTripper: authRoundTripper,
	})

	if err != nil {
		return nil, err
	}

	return &Client{client: v1.NewAPI(client), config: conf}, nil
}

func (c *Client) Query(ctx context.Context, query string) (model.Vector, error) {
	result, warnings, err := c.client.Query(ctx, query, time.Now(), v1.WithTimeout(time.Duration(c.config.Timeout)*time.Second))
	if err != nil {
		return nil, errors.Wrapf(err, "Error on prometheus query [%s]", query)
	}
	if len(warnings) > 0 {
		log.Log(ctx).Warnf("Warnings: %v", warnings)
	}

	s, ok := result.(model.Vector)
	if !ok {
		return nil, errors.Errorf("The query result must be Vector, query [%s]", query)
	}
	return s, nil
}

func (c *Client) QueryRange(ctx context.Context, query string, r v1.Range) (model.Matrix, error) {
	result, warnings, err := c.client.QueryRange(ctx, query, r, v1.WithTimeout(time.Duration(c.config.Timeout)*time.Second))
	if err != nil {
		return nil, errors.Wrapf(err, "Error on prometheus query [%s]", query)
	}
	if len(warnings) > 0 {
		log.Log(ctx).Warnf("Warnings: %v", warnings)
	}

	s, ok := result.(model.Matrix)
	if !ok {
		return nil, errors.Errorf("The query result must be Vector, query [%s]", query)
	}
	return s, nil
}
