package model

import (
	"fmt"
	"sort"
	"time"
)

// Resource represents the data model for a k8s resource.
type Resource struct {
	kind      string
	name      string
	createdAt time.Time
	restarts  int
}

// New returns a resource.
func New(kind string, name string, createdAt time.Time, restarts int) Resource {
	return Resource{
		kind:      kind,
		createdAt: createdAt,
		restarts:  restarts,
		name:      name,
	}
}

func (r Resource) String() string {
	return fmt.Sprintf("%s, %s, %d", r.name, r.createdAt, r.restarts)
}

func (r Resource) PrintResource() string {
	return fmt.Sprintf("%s, %s ago, %d", r.name, time.Since(r.createdAt), r.restarts)
}

type byAge []Resource

func (r byAge) Len() int           { return len(r) }
func (r byAge) Less(i, j int) bool { return r[i].createdAt.Before(r[j].createdAt) }
func (r byAge) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// SortByAge sorts all the resources as per their created timestamps
func SortByAge(rs []Resource) []Resource {
	sort.Sort(byAge(rs))
	return rs
}

type byName []Resource

func (r byName) Len() int           { return len(r) }
func (r byName) Less(i, j int) bool { return r[i].name < r[j].name }
func (r byName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// SortByName sorts all the resources as per their created timestamps
func SortByName(rs []Resource) []Resource {
	sort.Sort(byName(rs))
	return rs
}

type byRestart []Resource

func (r byRestart) Len() int           { return len(r) }
func (r byRestart) Less(i, j int) bool { return r[i].restarts < r[j].restarts }
func (r byRestart) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// SortByRestart sorts all the resources as per their created timestamps
func SortByRestart(rs []Resource) []Resource {
	sort.Sort(byRestart(rs))
	return rs
}
