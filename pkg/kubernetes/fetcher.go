package kubernetes

import (
	"context"

	"github.com/okteto/go-getting-started/pkg/kubernetes/client"
	"github.com/okteto/go-getting-started/pkg/kubernetes/fetcher"
	"github.com/okteto/go-getting-started/pkg/model"
)

type Fetcher interface {
	Fetch() []model.Resource
}

func GetAll(ctx context.Context) ([]model.Resource, error) {
	k8sClient, err := client.GetK8sRestclient()
	if err != nil {
		return nil, err
	}
	return fetcher.NewPods(ctx, k8sClient).Fetch()
}
