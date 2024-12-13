package workload

import (
	"context"
	"fmt"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"
)

func Run(ctx context.Context, runtime *runtime.Runtime, v model.ConnectorConfig) error {
	steps := *v.Steps
	switch steps.Kind() {
	case model.Inlet:
		return runInlet(ctx, runtime, steps)
	case model.Outlet:
		return runOutlet(ctx, runtime, steps)
	default:
		return fmt.Errorf("unknown connector kind %s", steps.Kind())
	}
}
