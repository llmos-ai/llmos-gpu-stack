# LLMOS-GPU-Stack Charts
![Release Charts](https://github.com/llmos-ai/llmos-gpu-stack/workflows/release/badge.svg)

## Prerequisites
- Helm 3.x

## Add Helm Repo
```shell
$ helm repo add llmos-gpu-stack https://llmos-gpu-stack-charts.1block.ai
$ helm repo update
```

## Installing the Charts
```shell
## Install CRDs
$ helm upgrade --install --create-namespace -n llmos-system llmos-gpu-stack-crds llmos-gpu-stack/llmos-gpu-stack-crds

## Install the Chart
$ helm upgrade --install --create-namespace -n llmos-system llmos-gpu-stack llmos-gpu-stack/llmos-gpu-stack --reuse-values
```