// Derived from https://raw.githubusercontent.com/biogo/store/master/llrb/llrb.go
//
// Copyright ©2012 The bíogo Authors. All rights reserved.
// Copyright ©2016 Markus Sonderegger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

// Element is a type that can be inserted into a Tree or used as a range
// or equality query on the tree.
type Element interface {
	// Compare returns a value indicating the sort order relationship
	// between the receiver and the parameter.
	//
	// Given c = a.Compare(b):
	//  c < 0 if a < b;
	//  c == 0 if a == b; and
	//  c > 0 if a > b.
	Compare(Element) int
}

// A Visitor is a function that operates on an Element. If done is
// returned true, the Visitor is indicating that no further work needs
// to be done and so the traversal function should traverse no further.
type Visitor func(Element) (done bool)

const (
	red   = false
	black = true
)

type node struct {
	elem  Element
	right *node
	left  *node
	color bool
}

func (n *node) copy() *node {
	return &node{
		elem:  n.elem,
		left:  n.left,
		right: n.right,
		color: n.color,
	}
}

func (n *node) rotateLeft() *node {
	root := n.right
	n.right = root.left
	root.left = n
	root.color = n.color
	n.color = red
	return root
}

func (n *node) rotateRight() *node {
	root := n.left
	n.left = root.right
	root.right = n
	root.color = n.color
	n.color = red
	return root
}

func (n *node) flipColors() {
	n.color = !n.color
	n.left.color = !n.left.color
	n.right.color = !n.right.color
}

func (n *node) isRed() bool {
	if n == nil {
		return false
	}
	return n.color == red
}

func (n *node) fixUp() *node {
	if n.right.isRed() {
		n = n.rotateLeft()
	}
	if n.left.isRed() && n.left.left.isRed() {
		n = n.rotateRight()
	}
	if n.left.isRed() && n.right.isRed() {
		n.flipColors()
	}
	return n
}

func (n *node) moveRedLeft() *node {
	n.flipColors()
	if n.right.left.isRed() {
		n.right = n.right.rotateRight()
		n = n.rotateLeft()
		n.flipColors()
	}
	return n
}

func (n *node) moveRedRight() *node {
	n.flipColors()
	if n.left.left.isRed() {
		n = n.rotateRight()
		n.flipColors()
	}
	return n
}

func (n *node) find(elem Element) *node {
	for n != nil {
		switch cmp := elem.Compare(n.elem); {
		case cmp == 0:
			return n
		case cmp < 0:
			n = n.left
		default:
			n = n.right
		}
	}
	return n
}

func (n *node) insert(elem Element) (*node, int) {
	if n == nil {
		return &node{elem: elem}, 1
	} else if n.elem == nil {
		n.elem = elem
		return n, 1
	}

	root, m := n.copy(), 0 // recursive branch copy
	switch cmp := elem.Compare(root.elem); {
	case cmp == 0:
		root.elem = elem
	case cmp < 0:
		root.left, m = root.left.insert(elem)
	default:
		root.right, m = root.right.insert(elem)
	}

	if root.right.isRed() && !root.left.isRed() {
		root = root.rotateLeft()
	}
	if root.left.isRed() && root.left.left.isRed() {
		root = root.rotateRight()
	}
	if root.left.isRed() && root.right.isRed() {
		root.flipColors()
	}
	return root, m
}

func (n *node) deleteMin() (*node, int) {
	if n.left == nil {
		return nil, -1
	}
	if !n.left.isRed() && !n.left.left.isRed() {
		n = n.moveRedLeft()
	}
	var m int
	n.left, m = n.left.deleteMin()

	root := n.fixUp()
	return root, m
}

func (n *node) deleteMax() (*node, int) {
	if n.left != nil && n.left.isRed() {
		n = n.rotateRight()
	}
	if n.right == nil {
		return nil, -1
	}
	if !n.right.isRed() && !n.right.left.isRed() {
		n = n.moveRedRight()
	}
	var m int
	n.right, m = n.right.deleteMax()

	root := n.fixUp()
	return root, m
}

func (n *node) min() *node {
	for ; n.left != nil; n = n.left {
	}
	return n
}

func (n *node) max() *node {
	for ; n.right != nil; n = n.right {
	}
	return n
}

func (n *node) delete(elem Element) (*node, int) {
	root, m := n.copy(), 0 // recursive branch copy

	if elem.Compare(root.elem) < 0 {
		if root.left != nil {
			if !root.left.isRed() && !root.left.left.isRed() {
				root = root.moveRedLeft()
			}
			root.left, m = root.left.delete(elem)
		}
	} else {
		if root.left.isRed() {
			root = root.rotateRight()
		}
		if root.right == nil && elem.Compare(root.elem) == 0 {
			return nil, -1
		}
		if root.right != nil {
			if !root.right.isRed() && !root.right.left.isRed() {
				root = root.moveRedRight()
			}
			if elem.Compare(root.elem) == 0 {
				root.elem = root.right.min().elem
				root.right, m = root.right.deleteMin()
			} else {
				root.right, m = root.right.delete(elem)
			}
		}
	}

	root = root.fixUp()
	return root, m
}

func (n *node) do(fn Visitor) (done bool) {
	if n.left != nil {
		done = n.left.do(fn)
		if done {
			return done
		}
	}
	if done = fn(n.elem); done {
		return done
	}
	if n.right != nil {
		done = n.right.do(fn)
	}
	return done
}

func (n *node) doRange(lo, hi Element, fn Visitor) (done bool) {
	lc, hc := lo.Compare(n.elem), hi.Compare(n.elem)
	if lc <= 0 && n.left != nil {
		done = n.left.doRange(lo, hi, fn)
		if done {
			return done
		}
	}
	if lc <= 0 && hc > 0 {
		if done = fn(n.elem); done {
			return
		}
	}
	if hc > 0 && n.right != nil {
		done = n.right.doRange(lo, hi, fn)
	}
	return done
}
