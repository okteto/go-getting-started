package model

import (
	"fmt"
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

func (r Resource) Pretty() string {
	return fmt.Sprintf("%s, %s ago, %d", r.name, time.Since(r.createdAt), r.restarts)
}

// ByAge sorts all the resources as per their created timestamps
type ByAge []Resource

func (r ByAge) Len() int           { return len(r) }
func (r ByAge) Less(i, j int) bool { return r[i].createdAt.Before(r[j].createdAt) }
func (r ByAge) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// ByName sorts all the resources as per their created timestamps
type ByName []Resource

func (r ByName) Len() int           { return len(r) }
func (r ByName) Less(i, j int) bool { return r[i].name < r[j].name }
func (r ByName) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

// ByRestarts sorts all the resources as per their created timestamps
type ByRestart []Resource

func (r ByRestart) Len() int           { return len(r) }
func (r ByRestart) Less(i, j int) bool { return r[i].restarts < r[j].restarts }
func (r ByRestart) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
