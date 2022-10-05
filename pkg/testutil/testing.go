package testutil

import (
	"sort"
	"testing"

	"github.com/okteto/go-getting-started/pkg/model"
)

func CheckErrorAndResources(t *testing.T, shouldErr bool, err error, expected, actual []model.Resource) {
	t.Helper()
	if err == nil && shouldErr {
		t.Error("expected error, but returned none")
		return
	}
	if err != nil && !shouldErr {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if len(actual) != len(expected) {
		t.Fatalf("length of the slices differ: Expected %d, but was %d", len(expected), len(actual))
	}
	sort.Sort(model.ByName(actual))
	sort.Sort(model.ByName(expected))
	for i, e := range expected {
		if actual[i].String() != e.String() {
			t.Fatalf("list differs: Expected %s, but was %s", e, actual[i])
		}
	}
}
