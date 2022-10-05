package render

import (
	"fmt"
	"net/http"

	"github.com/okteto/go-getting-started/pkg/kubernetes"
	"github.com/okteto/go-getting-started/pkg/model"
)

const (
	header = "Name, Created At, Restart Count"
)

// renderer is an empty struct right now.
// In future it can hold render config like namespace or kind objects
type renderer struct {
}

func New() renderer {
	return renderer{}
}

func (r renderer) SortByAge(w http.ResponseWriter, req *http.Request) {
	r.getSortedList(w, req, model.SortByAge)
}

func (r renderer) SortByRestart(w http.ResponseWriter, req *http.Request) {
	r.getSortedList(w, req, model.SortByRestart)
}

func (r renderer) SortByName(w http.ResponseWriter, req *http.Request) {
	r.getSortedList(w, req, model.SortByName)
}

func (r renderer) getSortedList(w http.ResponseWriter, req *http.Request, f func([]model.Resource) []model.Resource) {
	rs, err := kubernetes.GetAll(req.Context())
	if err != nil {
		// This method should check for different error types and decide appropriate http error code.
		http.Error(w, err.Error(), 500)
		return
	}
	if len(rs) == 0 {
		fmt.Fprintln(w, "No resources found")
		return
	}
	fmt.Fprintln(w, header)
	f(rs)
	for _, resource := range rs {
		fmt.Fprintln(w, resource.PrintResource())
	}
}
