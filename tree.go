// Derived from https://raw.githubusercontent.com/biogo/store/master/llrb/llrb.go
//
// Copyright ©2012 The bíogo Authors. All rights reserved.
// Copyright ©2016 Markus Sonderegger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package llrb implements an immutable Left-Leaning Red-Black tree as
// described by Robert Sedgewick. More details relating to the
// implementation are available at the following locations:
//
//  http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf
//  http://www.cs.princeton.edu/~rs/talks/LLRB/Java/RedBlackBST.java
//  http://www.teachsolaisgames.com/articles/balanced_left_leaning.html
//
// The immutable version of the llrb tree is obviously going to be slower
// than the mutable version but should offer higher read availability.
// Immutability is achieved by branch copying.
package llrb

// Tree manages the root node of an left-Leaning Red-Black  tree. Public
// methods are exposed through this type.
type Tree struct {
	root *node
	size int
}

// Txn is a transaction on the tree. This transaction is applied
// atomically and returns a new tree when committed. A transaction is not
// thread safe, and should only be used by a single goroutine.
type Txn struct {
	tree *Tree
}

// Range performs fn on all values stored in the tree over the interval
// [from, to) from left to right. If to is less than from Range will
// panic. A boolean is returned indicating whether the Range traversal
// was interrupted by an Visitor returning true. If fn alters stored
// values sort relationships future tree operation behaviors are
// undefined.
func (t *Tree) Range(from, to Element, fn Visitor) bool {
	if t.root == nil {
		return false
	}
	if from.Compare(to) > 0 {
		panic("inverted range")
	}
	return t.root.doRange(from, to, fn)
}

// ForEach performs fn on all values stored in the tree. A boolean is
// returned indicating whether the ForEach traversal was interrupted by
// a Visitor returning true. If fn alters stored values sort
// relationships, future tree operation behaviors are undefined.
func (t *Tree) ForEach(fn Visitor) bool {
	if t.root == nil {
		return false
	}
	return t.root.do(fn)
}

// Get returns the first match of elem in the Tree. If insertion without
// replacement is used, this is probably not what you want.
func (t *Tree) Get(elem Element) Element {
	if t.root == nil {
		return nil
	}
	n := t.root.find(elem)
	if n == nil {
		return nil
	}
	return n.elem
}

// Max returns the maximum value stored in the tree. This will be the
// right-most maximum value if insertion without replacement has been
// used.
func (t *Tree) Max() Element {
	if t.root == nil {
		return nil
	}
	return t.root.max().elem
}

// Min returns the minimum value stored in the tree. This will be the
// left-most minimum value if insertion without replacement has been
// used.
func (t *Tree) Min() Element {
	if t.root == nil {
		return nil
	}
	return t.root.min().elem
}

// Len returns the number of elements stored in the Tree.
func (t *Tree) Len() int { return t.size }

// Snapshot returns a copy of the underlying tree.
func (t *Tree) Snapshot() *Tree {
	tree := &Tree{}
	if t == nil {
		return tree
	}

	tree.size = t.size
	if t.root != nil {
		tree.root = t.root.copy()
	}
	return tree
}

// Txn starts a new transaction that can be used to mutate the tree.
func (t *Tree) Txn() *Txn {
	return &Txn{tree: t.Snapshot()}
}

// Commit is used to finalize the transaction and return a new tree
func (t *Txn) Commit() *Tree {
	return t.tree
}

// Get returns the first match of elem in the Tree. If insertion without
// replacement is used, this is probably not what you want.
func (t *Txn) Get(elem Element) Element {
	return t.tree.Get(elem)
}

// Max returns the maximum value stored in the tree. This will be the
// right-most maximum value if insertion without replacement has been
// used.
func (t *Txn) Max() Element {
	return t.tree.Max()
}

// Min returns the minimum value stored in the tree. This will be the
// left-most minimum value if insertion without replacement has been
// used.
func (t *Txn) Min() Element {
	return t.tree.Min()
}

// Insert inserts the Element elem into the Tree at the first match
// found with elem or when a nil node is reached. Insertion without
// replacement can specified by ensuring that elem.Compare() never
// returns 0. If insert without replacement is performed, a distinct
// query Element must be used that can return 0 with a elem.Compare()
// call.
func (t *Txn) Insert(elem Element) {
	root, m := t.tree.root.insert(elem)
	t.tree.size += m
	t.tree.root = root
	t.tree.root.color = black
}

// Delete deletes the node that matches elem according to Compare().
// Note that Compare must identify the target node uniquely and in cases
// where non-unique keys are used, attributes used to break ties must be
// used to determine tree ordering during insertion.
func (t *Txn) Delete(elem Element) {
	if t.tree == nil || t.tree.root == nil {
		return
	}
	root, m := t.tree.root.delete(elem)
	t.tree.size += m
	t.tree.root = root
	if root == nil {
		return
	}
	t.tree.root.color = black
}

// DeleteMax deletes the node with the maximum value in the tree. If
// insertion without replacement has been used, the right-most maximum
// will be deleted.
func (t *Txn) DeleteMax() {
	if t.tree == nil || t.tree.root == nil {
		return
	}
	root, m := t.tree.root.deleteMax()
	t.tree.size += m
	t.tree.root = root
	if root == nil {
		return
	}
	t.tree.root.color = black
}

// DeleteMin deletes the node with the minimum value in the tree. If
// insertion without replacement has been used, the left-most minimum
// will be deleted.
func (t *Txn) DeleteMin() {
	if t.tree == nil || t.tree.root == nil {
		return
	}
	root, m := t.tree.root.deleteMin()
	t.tree.size += m
	t.tree.root = root
	if root == nil {
		return
	}
	t.tree.root.color = black
}

// Len returns the number of elements stored in the Tree.
func (t *Txn) Len() int { return t.tree.size }
