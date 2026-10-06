package main

import (
	"container/heap"
)

// Node represents a node in the Huffman tree.
type Node struct {
	freq  int
	leaf  bool
	val   byte
	left  *Node
	right *Node
}

// priorityQueue implements heap.Interface for *Node based on freq.
type priorityQueue []*Node

func (pq priorityQueue) Len() int { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].freq < pq[j].freq
}
func (pq priorityQueue) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }

func (pq *priorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// BuildTree constructs a Huffman tree from the given frequency map.
func BuildTree(freq map[byte]int) *Node {
	pq := &priorityQueue{}
	heap.Init(pq)
	for b, f := range freq {
		if f > 0 {
			heap.Push(pq, &Node{freq: f, leaf: true, val: b})
		}
	}
	if pq.Len() == 0 {
		return nil
	}
	for pq.Len() > 1 {
		n1 := heap.Pop(pq).(*Node)
		n2 := heap.Pop(pq).(*Node)
		parent := &Node{freq: n1.freq + n2.freq, leaf: false, left: n1, right: n2}
		heap.Push(pq, parent)
	}
	return heap.Pop(pq).(*Node)
}

// buildCodeMap traverses the tree to build a map from byte to its bit string.
func buildCodeMap(node *Node, prefix string, m map[byte]string) {
	if node == nil {
		return
	}
	if node.leaf {
		m[node.val] = prefix
		return
	}
	buildCodeMap(node.left, prefix+"0", m)
	buildCodeMap(node.right, prefix+"1", m)
}

// Encode encodes data using Huffman coding. Returns the bit string and the code map.
func Encode(data []byte, root *Node) (string, map[byte]string) {
	if root == nil {
		return "", nil
	}
	codeMap := make(map[byte]string)
	buildCodeMap(root, "", codeMap)
	var encoded string
	for _, b := range data {
		encoded += codeMap[b]
	}
	return encoded, codeMap
}

// Decode decodes a bit string using the Huffman tree.
func Decode(encoded string, root *Node) []byte {
	if root == nil || encoded == "" {
		return []byte{}
	}
	var decoded []byte
	node := root
	for _, bit := range encoded {
		if bit == '0' {
			node = node.left
		} else {
			node = node.right
		}
		if node.leaf {
			decoded = append(decoded, node.val)
			node = root
		}
	}
	return decoded
}
