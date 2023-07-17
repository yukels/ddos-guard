package prometheus

import (
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"

	"github.com/yukels/util/context"
)

type Provider interface {
	Query(ctx context.Context, query string) (model.Vector, error)
	QueryRange(ctx context.Context, query string, r v1.Range) (model.Matrix, error)
}
