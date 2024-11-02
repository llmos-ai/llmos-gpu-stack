package devicemanager

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/llmos-ai/llmos/utils/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/llmos-ai/llmos-gpu-stack/pkg/controller"
)

type DeviceManager struct {
	Kubeconfig      string `usage:"Path to kubeconfig file" short:"k" env:"KUBECONFIG"`
	Namespace       string `usage:"Deployment namespace" default:"llmos-system" short:"n" env:"LLMOS_NAMESPACE"`
	HttpAddress     string `usage:"Address to listen for HTTP requests" default:"0.0.0.0:8080" short:"a" env:"LLMOS_HTTP_ADDRESS"`
	ProfilerAddress string `usage:"Address to listen for profiling" default:"0.0.0.0:6060" short:"p" env:"LLMOS_PROFILER_ADDRESS"`
	Threadiness     int    `usage:"Number of threads to run" default:"2" short:"w" env:"LLMOS_THREADINESS"`
}

func NewManager() *cobra.Command {
	return cli.Command(&DeviceManager{}, cobra.Command{
		Short: "Run device manager",
	})
}

func (dm *DeviceManager) Run(cmd *cobra.Command, _ []string) error {
	initProfiling(dm)
	gc := controller.GPUDeviceController{
		Context:     cmd.Context(),
		Kubeconfig:  dm.Kubeconfig,
		HttpAddress: dm.HttpAddress,
		Threadiness: dm.Threadiness,
	}
	return gc.Start()
}

func initProfiling(dm *DeviceManager) {
	// enable profiler
	if dm.ProfilerAddress != "" {
		logrus.Infof("starting profiler server on %s", dm.ProfilerAddress)
		profilerServer := &http.Server{
			Addr:              dm.ProfilerAddress,
			ReadHeaderTimeout: 10 * time.Second,
		}
		go func() {
			log.Println(profilerServer.ListenAndServe())
		}()
	}
}
