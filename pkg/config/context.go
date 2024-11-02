package config

import (
	"context"

	hdevice "github.com/Project-HAMi/HAMi/pkg/device"
	"github.com/rancher/lasso/pkg/controller"
	appsv1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/apps"
	corev1 "github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	"github.com/rancher/wrangler/v3/pkg/generic"
	"github.com/rancher/wrangler/v3/pkg/ratelimit"
	"github.com/rancher/wrangler/v3/pkg/start"
	"k8s.io/client-go/kubernetes"

	gpustackv1 "github.com/llmos-ai/llmos-gpu-stack/pkg/generated/controllers/gpustack.llmos.ai"
)

type Management struct {
	Ctx        context.Context
	K8s        *kubernetes.Clientset
	GPUDevices map[string]hdevice.Devices

	CoreFactory     *corev1.Factory
	AppsFactory     *appsv1.Factory
	GPUStackFactory *gpustackv1.Factory

	starters []start.Starter
}

func NewManagementContext(ctx context.Context, kubeConfig string) (*Management, error) {
	hdevice.InitDevices()
	hdevice.GlobalFlagSet()
	mgmt := &Management{
		Ctx:        ctx,
		GPUDevices: hdevice.GetDevices(),
	}

	clientConfig, err := GetConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	client, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	client.RateLimiter = ratelimit.None

	k8s, err := kubernetes.NewForConfig(client)
	if err != nil {
		return nil, err
	}
	mgmt.K8s = k8s

	factory, err := controller.NewSharedControllerFactoryFromConfig(client, nil)
	if err != nil {
		return nil, err
	}

	opts := &generic.FactoryOptions{
		SharedControllerFactory: factory,
	}

	core, err := corev1.NewFactoryFromConfigWithOptions(client, opts)
	if err != nil {
		return nil, err
	}
	mgmt.CoreFactory = core
	mgmt.starters = append(mgmt.starters, core)

	apps, err := appsv1.NewFactoryFromConfigWithOptions(client, opts)
	if err != nil {
		return nil, err
	}
	mgmt.AppsFactory = apps
	mgmt.starters = append(mgmt.starters, apps)

	gpuStack, err := gpustackv1.NewFactoryFromConfigWithOptions(client, opts)
	if err != nil {
		return nil, err
	}
	mgmt.GPUStackFactory = gpuStack
	mgmt.starters = append(mgmt.starters, gpuStack)

	return mgmt, nil
}

func (m *Management) Start(threads int) error {
	return start.All(m.Ctx, threads, m.starters...)
}
