package main

import (
	"errors"
	"strings"
)

type Node struct {
	name     string
	isDir    bool
	children map[string]*Node
}

func NewNode(name string, isDir bool) *Node {
	return &Node{
		name:  name,
		isDir: isDir,
	}
}

func (n *Node) AddPath(path string) error {
	parts := strings.Split(path, "/")
	current := n
	for i, part := range parts {
		if part == "" {
			continue
		}
		child, exists := current.children[part]
		if !exists {
			isDir := i < len(parts)-1
			child = NewNode(part, isDir)
			if current.children == nil {
				current.children = make(map[string]*Node)
			}
			current.children[part] = child
		}
		current = child
	}
	return nil
}

func (n *Node) Find(path string) (*Node, error) {
	if path == "" || path == "/" {
		return n, nil
	}
	parts := strings.Split(path, "/")
	current := n
	for _, part := range parts {
		if part == "" {
			continue
		}
		child, exists := current.children[part]
		if !exists {
			return nil, errors.New("path not found")
		}
		current = child
	}
	return current, nil
}

func (n *Node) List() []string {
	if n.children == nil {
		return nil
	}
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	return names
}
