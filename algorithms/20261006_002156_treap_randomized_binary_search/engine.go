package main

// merge combines two treaps where all keys in left are less than those in right.
func merge(a, b *Node) *Node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.priority > b.priority {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

// split divides the treap rooted at node into two treaps: keys < key and keys >= key.
func split(node *Node, key int) (*Node, *Node) {
	if node == nil {
		return nil, nil
	}
	if key <= node.key {
		left, right := split(node.left, key)
		node.left = right
		return left, node
	}
	left, right := split(node.right, key)
	node.right = left
	return node, right
}

// deleteNode removes a key from the treap rooted at node.
func deleteNode(node *Node, key int) *Node {
	if node == nil {
		return nil
	}
	if key < node.key {
		node.left = deleteNode(node.left, key)
	} else if key > node.key {
		node.right = deleteNode(node.right, key)
	} else {
		return merge(node.left, node.right)
	}
	return node
}

// searchNode checks whether key exists in the treap rooted at node.
func searchNode(node *Node, key int) bool {
	if node == nil {
		return false
	}
	if key == node.key {
		return true
	}
	if key < node.key {
		return searchNode(node.left, key)
	}
	return searchNode(node.right, key)
}

// inOrder traverses the treap in-order and appends keys to res.
func inOrder(node *Node, res *[]int) {
	if node == nil {
		return
	}
	inOrder(node.left, res)
	*res = append(*res, node.key)
	inOrder(node.right, res)
}
