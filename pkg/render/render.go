package render

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/okteto/go-getting-started/pkg/kubernetes"
	"github.com/okteto/go-getting-started/pkg/model"
)

// renderer is an empty struct right now.
// In future it can hold render config like namespace or kind objects
type renderer struct {
}

func New() renderer {
	return renderer{}

}

func (r renderer) SortByAge(w http.ResponseWriter, req *http.Request) {
	rs, err := kubernetes.GetAll(req.Context())
	if err != nil {
		// This method should check for different error types and decide appropriate http error code.
		http.Error(w, err.Error(), 500)
		return
	}
	sort.Sort(model.ByAge(rs))
	for _, resource := range rs {
		fmt.Fprint(w, resource)

	}
}

func (r renderer) SortByRestart(w http.ResponseWriter, req *http.Request) {
	rs, err := kubernetes.GetAll(req.Context())
	if err != nil {
		// This method should check for different error types and decide appropriate http error code.
		http.Error(w, err.Error(), 500)
		return
	}
	sort.Sort(model.ByRestart(rs))
	for _, resource := range rs {
		fmt.Fprint(w, resource)
	}
}

func (r renderer) SortByName(w http.ResponseWriter, req *http.Request) {
	rs, err := kubernetes.GetAll(req.Context())
	if err != nil {
		// This method should check for different error types and decide appropriate http error code.
		http.Error(w, err.Error(), 500)
		return
	}
	sort.Sort(model.ByName(rs))
	for _, resource := range rs {
		fmt.Fprint(w, resource)
	}
}
