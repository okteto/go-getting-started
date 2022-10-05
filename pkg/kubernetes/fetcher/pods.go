package fetcher

import (
	"context"

	"github.com/okteto/go-getting-started/pkg/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	podKind = "Pod"
)

func NewPods(ctx context.Context, restClient kubernetes.Interface) podsFetcher {
	return podsFetcher{
		ctx:        ctx,
		restClient: restClient,
	}
}

type podsFetcher struct {
	ctx        context.Context
	restClient kubernetes.Interface
}

func (p podsFetcher) Fetch() ([]model.Resource, error) {

	// namespace should be read from current context
	// TODO (tejal29): read namespace.
	pods, err := p.restClient.CoreV1().Pods("").List(p.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	results := make([]model.Resource, len(pods.Items))
	for i, pod := range pods.Items {
		createdAt := pod.ObjectMeta.CreationTimestamp.Time
		restartCnt := 0
		// Go through all init containers and container statuses to find the restart count
		for _, c := range append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...) {
			restartCnt = restartCnt + int(c.RestartCount)
		}
		results[i] = model.New(podKind, pod.GetName(), createdAt, restartCnt)

	}
	return results, nil
}
