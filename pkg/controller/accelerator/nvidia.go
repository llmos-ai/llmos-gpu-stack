package accelerator

import corev1 "k8s.io/api/core/v1"

const (
	gpuPresentLabel = "nvidia.com/gpu.present"
)

func NodeHasGPUPresent(node *corev1.Node) bool {
	return node.Labels[gpuPresentLabel] == "true"
}
