package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/llmos-ai/llmos-gpu-stack/cmd"
)

func main() {
	cmd := cmd.NewRootCmd()
	ctx := signals.SetupSignalHandler()
	cmd.SilenceErrors = true
	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Errorf("%s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
