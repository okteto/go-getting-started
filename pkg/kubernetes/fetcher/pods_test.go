package fetcher

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"

	"github.com/okteto/go-getting-started/pkg/model"
	"github.com/okteto/go-getting-started/pkg/testutil"
)

func TestPodFetch(t *testing.T) {
	now := time.Now()
	anHourAgo := now.Add(time.Hour * -1)
	aMinAgo := now.Add(time.Minute * -1)
	tests := []struct {
		description string
		pods        []v1.Pod
		expected    []model.Resource
		shouldErr   bool
	}{
		{
			description: "multiple pods with containers running",
			pods: []v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo", CreationTimestamp: metav1.Time{Time: aMinAgo}},
					Status: v1.PodStatus{
						ContainerStatuses: []v1.ContainerStatus{{Ready: true, RestartCount: 2}},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{Name: "bar", CreationTimestamp: metav1.Time{Time: anHourAgo}},
					Status: v1.PodStatus{
						ContainerStatuses: []v1.ContainerStatus{{Ready: true, RestartCount: 1}},
					},
				},
			},
			expected: []model.Resource{
				model.New(podKind, "foo", aMinAgo, 2),
				model.New(podKind, "bar", anHourAgo, 1),
			},
		},
		{
			description: "no pods",
			pods:        []v1.Pod{},
			expected:    []model.Resource{},
		},
		{
			description: "pods with multiple container status",
			pods: []v1.Pod{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "foo", CreationTimestamp: metav1.Time{Time: aMinAgo}},
					Status: v1.PodStatus{
						InitContainerStatuses: []v1.ContainerStatus{{Ready: false, RestartCount: 1}, {Ready: false, RestartCount: 1}},
						ContainerStatuses:     []v1.ContainerStatus{{Ready: false, RestartCount: 2}, {Ready: false, RestartCount: 1}},
					},
				},
			},
			expected: []model.Resource{model.New(podKind, "foo", aMinAgo, 5)},
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			objs := make([]runtime.Object, len(test.pods))
			for i, pod := range test.pods {
				objs[i] = pod.DeepCopy()
			}
			client := fakekubeclientset.NewSimpleClientset(objs...)
			actual, err := NewPods(context.Background(), client).Fetch()

			testutil.CheckErrorAndResources(t, test.shouldErr, err, test.expected, actual)
		})
	}
}
