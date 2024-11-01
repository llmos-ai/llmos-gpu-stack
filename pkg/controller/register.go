package controller

import (
	"context"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/config"
	"github.com/llmos-ai/llmos-gpu-stack/pkg/controller/gpudevice"
)

type registerFunc func(context.Context, *config.Management) error

var registerFuncs = []registerFunc{
	gpudevice.Register,
}

func register(ctx context.Context, mgmt *config.Management) error {
	for _, f := range registerFuncs {
		if err := f(ctx, mgmt); err != nil {
			return err
		}
	}
	return nil
}
