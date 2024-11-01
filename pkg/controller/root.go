package controller

import (
	"context"

	"github.com/rancher/wrangler/v3/pkg/leader"
	"github.com/sirupsen/logrus"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/config"
)

const gpuStackLeaseName = "llmos-gpu-stack-device-manager-leader"

func Start(ctx context.Context, kubeConfig string, threadiness int) error {
	mgmt, err := config.NewManagementContext(ctx, kubeConfig)
	if err != nil {
		return err
	}

	leader.RunOrDie(ctx, "", gpuStackLeaseName, mgmt.K8s, func(ctx context.Context) {
		if err = register(ctx, mgmt); err != nil {
			logrus.Fatal(err)
		}
		if err = mgmt.Start(threadiness); err != nil {
			logrus.Fatal(err)
		}
	})

	<-ctx.Done()
	return ctx.Err()
}
