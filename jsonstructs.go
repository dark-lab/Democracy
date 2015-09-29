package main

type Graph struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

type Node struct {
	Name      string  `json:"name"`
	Group     int     `json:"group"`
	Size      float32 `json:"size"`
	Thickness float32 `json:"thickness"`
}

type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}
