package prometheus

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

// Handler for process metrics
type Handler struct {
	handler  http.Handler
	registry *prometheus.Registry
}

// NewHandler constructor
func NewHandler(ctx context.Context) (*Handler, error) {
	h := &Handler{
		registry: Registry(),
	}
	innerHandler, err := h.innerHandler()
	if err != nil {
		return nil, errors.Wrapf(err, "Couldn't create metrics handler")
	}
	h.handler = innerHandler
	return h, nil
}

// Metrics implements http.Handler.
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r.Context())

	log.Log(ctx).Debug("collect query")

	h.handler.ServeHTTP(w, r)
}

// innerHandler is used to create both the one unfiltered http.Handler to be
// wrapped by the outer handler and also the filtered handlers created on the
// fly. The former is accomplished by calling innerHandler without any arguments
func (h *Handler) innerHandler() (http.Handler, error) {
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{h.registry},
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
			Registry:      h.registry,
		},
	)
	// Note that we have to use h.registry here to
	// use the same promhttp metrics for all expositions.
	handler = promhttp.InstrumentMetricHandler(
		h.registry, handler,
	)
	return handler, nil
}
