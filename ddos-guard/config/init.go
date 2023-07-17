package config

import (
	configutil "github.com/yukels/util/config"
	"github.com/yukels/util/context"
)

func init() {
	ctx := context.Background()
	configutil.InitFromDefaultDir(ctx, Sections)
}
