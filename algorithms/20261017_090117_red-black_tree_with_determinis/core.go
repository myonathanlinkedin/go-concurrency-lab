package main

type node struct {
	key   int
	val   string
	left  *node
	right *node
	color bool // true = red, false = black
	size  int
}

type RBTree struct {
	root *node
}

const (
	red   = true
	black = false
)

// isRed returns true if the node is non-nil and red.
func isRed(x *node) bool {
	if x == nil {
		return false
	}
	return x.color == red
}

// size returns the number of nodes in the subtree rooted at x.
func size(x *node) int {
	if x == nil {
		return 0
	}
	return x.size
}

// rotateLeft performs a left rotation.
func rotateLeft(h *node) *node {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = red
	x.size = h.size
	h.size = 1 + size(h.left) + size(h.right)
	return x
}

// rotateRight performs a right rotation.
func rotateRight(h *node) *node {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = red
	x.size = h.size
	h.size = 1 + size(h.left) + size(h.right)
	return x
}

// flipColors flips colors of a node and its two children.
func flipColors(h *node) {
	h.color = !h.color
	if h.left != nil {
		h.left.color = !h.left.color
	}
	if h.right != nil {
		h.right.color = !h.right.color
	}
}

// moveRedLeft prepares for deletion of the minimum key.
func moveRedLeft(h *node) *node {
	flipColors(h)
	if isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
		flipColors(h)
	}
	return h
}

// moveRedRight prepares for deletion of the maximum key.
func moveRedRight(h *node) *node {
	flipColors(h)
	if isRed(h.left.left) {
		h = rotateRight(h)
		flipColors(h)
	}
	return h
}

// fixUp restores LLRB invariants after modifications.
func fixUp(h *node) *node {
	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}
	h.size = 1 + size(h.left) + size(h.right)
	return h
}

// Insert adds a key/value pair to the tree, overwriting existing value if key exists.
func (t *RBTree) Insert(key int, val string) {
	t.root = put(t.root, key, val)
	t.root.color = black
}

// put recursively inserts a node.
func put(h *node, key int, val string) *node {
	if h == nil {
		return &node{key: key, val: val, color: red, size: 1}
	}
	if key < h.key {
		h.left = put(h.left, key, val)
	} else if key > h.key {
		h.right = put(h.right, key, val)
	} else {
		h.val = val
	}
	// restore LLRB properties
	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}
	h.size = 1 + size(h.left) + size(h.right)
	return h
}

// Get retrieves the value associated with the given key.
func (t *RBTree) Get(key int) (string, bool) {
	x := t.root
	for x != nil {
		if key < x.key {
			x = x.left
		} else if key > x.key {
			x = x.right
		} else {
			return x.val, true
		}
	}
	return "", false
}

// Delete removes the node with the given key, if present.
func (t *RBTree) Delete(key int) {
	if t.root == nil {
		return
	}
	if !isRed(t.root.left) && !isRed(t.root.right) {
		t.root.color = red
	}
	t.root = deleteNode(t.root, key)
	if t.root != nil {
		t.root.color = black
	}
}

// deleteNode recursively deletes a key.
func deleteNode(h *node, key int) *node {
	if key < h.key {
		if h.left != nil {
			if !isRed(h.left) && !isRed(h.left.left) {
				h = moveRedLeft(h)
			}
			h.left = deleteNode(h.left, key)
		}
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
		}
		if key == h.key && h.right == nil {
			return nil
		}
		if h.right != nil {
			if !isRed(h.right) && !isRed(h.right.left) {
				h = moveRedRight(h)
			}
			if key == h.key {
				minNode := min(h.right)
				h.key = minNode.key
				h.val = minNode.val
				h.right = deleteMin(h.right)
			} else {
				h.right = deleteNode(h.right, key)
			}
		}
	}
	return fixUp(h)
}

// min returns the node with the minimum key in the subtree.
func min(h *node) *node {
	for h.left != nil {
		h = h.left
	}
	return h
}

// deleteMin removes the minimum key node.
func deleteMin(h *node) *node {
	if h.left == nil {
		return nil
	}
	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}
	h.left = deleteMin(h.left)
	return fixUp(h)
}

// Validate checks all red‑black invariants.
// Returns true if the tree satisfies:
// 1) Root is black.
// 2) No red node has a red child.
// 3) Every path from root to leaf has the same number of black nodes.
func (t *RBTree) Validate() bool {
	if t.root != nil && t.root.color != black {
		return false
	}
	_, ok := validateNode(t.root)
	return ok
}

// validateNode returns (blackHeight, ok) for the subtree.
func validateNode(x *node) (int, bool) {
	if x == nil {
		return 1, true // null leaves count as black
	}
	if isRed(x) {
		if isRed(x.left) || isRed(x.right) {
			return 0, false // red violation
		}
	}
	leftBH, okLeft := validateNode(x.left)
	if !okLeft {
		return 0, false
	}
	rightBH, okRight := validateNode(x.right)
	if !okRight {
		return 0, false
	}
	if leftBH != rightBH {
		return 0, false // black‑height mismatch
	}
	if isRed(x) {
		return leftBH, true
	}
	return leftBH + 1, true
}
