# LLMOS-GPU-Stack
[![main-build](https://github.com/llmos-ai/llmos-gpu-stack/actions/workflows/main-release.yaml/badge.svg)](https://github.com/llmos-ai/llmos-gpu-stack/actions/workflows/main-release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/llmos-ai/llmos-gpu-stack)](https://goreportcard.com/report/github.com/llmos-ai/llmos-gpu-stack)
[![Releases](https://img.shields.io/github/release/llmos-ai/llmos-gpu-stack.svg)](https://github.com/llmos-ai/llmos-gpu-stack/releases)


LLMOS-GPU-Stack is a collection of tools that provides vGPU and Multi-accelerator support for the [LLMOS](https://github.com/llmos-ai/llmos) project.

## Getting Started

### Prerequisites
- Go version v1.22.0+
- Kubectl version v1.29.0+.
- Access to a Kubernetes v1.29.0+ cluster.
- Helm v3.0.0+

### Installation
To deploy the `llmos-gpu-stack` on your k8s cluster, you can use the following commands:

**Clone the Repo and install the llmos-gpu-stack & dependency charts to the cluster:**

```sh
$ make install
```

### Uninstall
**Delete the CRDs and llmos-gpu-stack from the cluster:**

```sh
$ make uninstall
```

## Helm Repo

If you want to use the Helm chart directly, you can add the repo and search for the chart:

```shell
helm repo add llmos-gpu-stack https://llmos-gpu-stack-charts.1block.ai
helm repo update llmos-gpu-stack
helm search repo llmos-gpu-stack # append `--devel` to list dev versions
```

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

