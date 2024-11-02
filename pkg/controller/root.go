package controller

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/rancher/wrangler/v3/pkg/leader"
	"github.com/sirupsen/logrus"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/config"
)

const gpuStackLeaseName = "llmos-gpu-stack-device-manager-leader"

type GPUDeviceController struct {
	Context     context.Context
	Kubeconfig  string
	HttpAddress string
	Threadiness int
}

func (c *GPUDeviceController) Start() error {
	mgmt, err := config.NewManagementContext(c.Context, c.Kubeconfig)
	if err != nil {
		return err
	}

	go leader.RunOrDie(c.Context, "", gpuStackLeaseName, mgmt.K8s, func(ctx context.Context) {
		if err = register(ctx, mgmt); err != nil {
			logrus.Fatal(err)
		}
		if err = mgmt.Start(c.Threadiness); err != nil {
			logrus.Fatal(err)
		}
		<-c.Context.Done()
	})

	// Start http server
	c.listenAndServe()

	<-c.Context.Done()
	return c.Context.Err()
}

func (c *GPUDeviceController) listenAndServe() {
	h := server.New(server.WithHostPorts(c.HttpAddress))

	h.GET("/healthz", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "ok"})
	})
	h.Spin()
}
