package main

import (
	"os"

	controllergen "github.com/rancher/wrangler/v3/pkg/controller-gen"
	"github.com/rancher/wrangler/v3/pkg/controller-gen/args"
)

func main() {
	_ = os.Unsetenv("GOPATH")
	controllergen.Run(args.Options{
		OutputPackage: "github.com/llmos-ai/llmos-gpu-stack/pkg/generated",
		Boilerplate:   "hack/boilerplate.go.txt",
		Groups: map[string]args.Group{
			"gpustack.llmos.ai": {
				PackageName: "gpustack.llmos.ai",
				Types: []interface{}{
					// All structs with an embedded ObjectMeta field will be picked up
					"./pkg/apis/gpustack.llmos.ai/v1",
				},
				GenerateTypes:   true,
				GenerateClients: true,
			},
		},
	})
}
