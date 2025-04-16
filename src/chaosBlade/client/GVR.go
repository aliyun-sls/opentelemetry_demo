package client

import "k8s.io/apimachinery/pkg/runtime/schema"

type Gvrpod struct {
	Gvr schema.GroupVersionResource
}

type Gvrnode struct {
	Gvr schema.GroupVersionResource
}

func (g *Gvrpod) GvrPodInit() {
	g.Gvr = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
}

func (g *Gvrnode) GvrNodeInit() {
	g.Gvr = schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}
}
