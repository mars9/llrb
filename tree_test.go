// Derived from https://raw.githubusercontent.com/biogo/store/master/llrb/llrb_test.go
//
// Copyright ©2012 The bíogo Authors. All rights reserved.
// Copyright ©2016 Markus Sonderegger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math/rand"
	"reflect"
	"testing"
)

func (t *Tree) is23() bool {
	if t == nil {
		return true
	}
	return t.root.is23()
}

func (t *Tree) isBalanced() bool {
	if t == nil {
		return true
	}
	var black int // number of black links on path from root to min
	for x := t.root; x != nil; x = x.left {
		if !x.isRed() {
			black++
		}
	}
	return t.root.isBalanced(black)
}

func (t *Tree) isBST() bool {
	if t == nil {
		return true
	}
	return t.root.isBST(t.Min(), t.Max())
}

func TestNilOperation(t *testing.T) {
	tree := &Tree{}
	if tree.Min() != nil {
		t.Fatalf("expected <nil> value, got %v", tree.Min())
	}
	if tree.Max() != nil {
		t.Fatalf("expected <nil> value, got %v", tree.Max())
	}

	txn := tree.Txn()
	txn.DeleteMin()
	if !reflect.DeepEqual(tree, txn.Commit()) {
		t.Fatalf("expected empty tree, got %#v", txn.Commit())
	}

	txn = tree.Txn()
	txn.DeleteMax()
	if !reflect.DeepEqual(tree, txn.Commit()) {
		t.Fatalf("expected empty tree, got %#v", txn.Commit())
	}

	if tree.Get(Int(42)) != nil {
		t.Fatalf("expected <nil> value, got %v", tree.Get(Int(42)))
	}
	if txn.Get(Int(42)) != nil {
		t.Fatalf("expected <nil> value, got %v", txn.Get(Int(42)))
	}

	txn.Delete(Int(42))
	if !reflect.DeepEqual(tree, txn.Commit()) {
		t.Fatalf("expected empty tree, got %#v", txn.Commit())
	}
}

func TestInsertion(t *testing.T) {
	min, max := compRune(0), compRune(1000)
	tree := &Tree{}
	txn := tree.Txn()
	for i := min; i <= max; i++ {
		txn.Insert(i)
		if txn.Len() != int(i+1) {
			t.Fatalf("insertion: expected tree length %d, have %d", i+1, txn.Len())
		}
		if !txn.tree.isBST() {
			t.Fatalf("insertion: tree is not a BST")
		}
		if !txn.tree.isBalanced() {
			t.Fatalf("insertion: tree is not balanced")
		}
		if !txn.tree.is23() {
			t.Fatalf("insertion: invariant violation")
		}
	}

	tree = txn.Commit()
	if tree.Min() != min {
		t.Fatalf("insertion: expected min element %d, have %d", min, tree.Min())
	}
	if tree.Max() != max {
		t.Fatalf("insertion: expected max element %d, have %d", min, tree.Max())
	}
}

func TestDeletion(t *testing.T) {
	min, max := compRune(0), compRune(1000)
	elems := int(max-min) + 1
	tree := &Tree{}
	txn := tree.Txn()
	for i := min; i <= max; i++ {
		txn.Insert(i)
	}
	tree = txn.Commit()

	txn = tree.Txn()
	for i := min; i <= max; i++ {
		if txn.Get(i) != nil {
			elems--
		}
		txn.Delete(i)
		if txn.Len() != elems {
			t.Fatalf("deletion: expected tree length %d, have %d", elems, txn.Len())
		}
		if i < max {
			if !txn.tree.isBST() {
				t.Fatalf("deletion: tree is not a BST")
			}
			if !txn.tree.isBalanced() {
				t.Fatalf("deletion: tree is not balanced")
			}
			if !txn.tree.is23() {
				t.Fatalf("deletion: invariant violation")
			}
		}
	}

	tree = txn.Commit()
	if tree.Len() != 0 {
		t.Fatalf("deletion: expected empty tree, have %d", tree.Len())
	}
}

func TestGet(t *testing.T) {
	min, max := compRune(0), compRune(1000)
	tree := &Tree{}
	txn := tree.Txn()
	for i := min; i <= max; i++ {
		if i&1 == 0 {
			txn.Insert(i)
		}
	}
	tree = txn.Commit()

	for i := min; i <= max; i++ {
		if i&1 == 0 {
			if tree.Get(i) != compRune(i) {
				t.Fatalf("get: expected element %v, got %v", compRune(i), tree.Get(i))
			}
		} else {
			if tree.Get(i) != nil {
				t.Fatalf("get: unexpected elem found %v", tree.Get(i))
			}
		}
	}
}

func TestRandomlyInsertedGet(t *testing.T) {
	count, max := 100000, 1000
	tree := &Tree{}
	txn := tree.Txn()
	verify := map[rune]struct{}{}
	for i := 0; i < count; i++ {
		v := compRune(rand.Intn(max))
		txn.Insert(v)
		verify[rune(v)] = struct{}{}
	}
	tree = txn.Commit()

	for v := range verify {
		if tree.Get(compRune(v)) != compRune(v) {
			t.Fatalf("random inserted: expected elem %v, got %v", compRune(v), tree.Get(compRune(v)))
		}
	}

	for i := compRune(0); i <= compRune(max); i++ {
		if _, ok := verify[rune(i)]; ok {
			if tree.Get(i) != i {
				t.Fatalf("random inserted: expected elem %v, got %v", i, tree.Get(i))
			}
		} else {
			if tree.Get(i) != nil {
				t.Errorf("get: unexpected elem found %v", tree.Get(i))
			}
		}
	}
}

func TestRandomInsertionAndDeltion(t *testing.T) {
	count, max := 100000, 1000
	tree := &Tree{}
	r := map[compRune]struct{}{}
	for i := 0; i < count; i++ {
		txn := tree.Txn()
		v := compRune(rand.Intn(max))
		r[v] = struct{}{}
		txn.Insert(v)
		tree = txn.Commit()

		if !tree.isBST() {
			t.Fatalf("random insertion and deletion: tree is not a BST")
		}
		if !tree.isBalanced() {
			t.Fatalf("random insertion and deletion: tree is not balanced")
		}
		if !tree.is23() {
			t.Fatalf("random insertion and deletion: tree is not a 2-3 tree")
		}
	}

	for v := range r {
		txn := tree.Txn()
		txn.Delete(v)
		tree = txn.Commit()

		if !tree.isBST() {
			t.Fatalf("random insertion and deletion: tree is not a BST")
		}
		if !tree.isBalanced() {
			t.Fatalf("random insertion and deletion: tree is not balanced")
		}
		if !tree.is23() {
			t.Fatalf("random insertion and deletion: tree is not a 2-3 tree")
		}
	}

	if tree.Len() != 0 {
		t.Fatalf("random insertion and deletion: expected empty tree, have %d", tree.Len())
	}
}

func TestDeleteMinMax(t *testing.T) {
	min, max := compRune(0), compRune(10)
	tree := &Tree{}
	txn := tree.Txn()
	for i := min; i <= max; i++ {
		txn.Insert(i)
	}
	tree = txn.Commit()
	if tree.Len() != int(max-min+1) {
		t.Fatalf("delete min/max: expected tree length %d, have %d", int(max-min+1), tree.Len())
	}

	for i, m := 0, int(max); i < m/2; i++ {
		txn = tree.Txn()
		txn.DeleteMin()
		tree = txn.Commit()

		if !tree.isBST() {
			t.Fatalf("delete min/max: tree is not a BST")
		}
		if !tree.isBalanced() {
			t.Fatalf("delete min/max: tree is not balanced")
		}
		if !tree.is23() {
			t.Fatalf("delete min/max: tree is not a 2-3 tree")
		}

		txn = tree.Txn()
		txn.DeleteMax()
		tree = txn.Commit()

		if !tree.isBST() {
			t.Fatalf("delete min/max: tree is not a BST")
		}
		if !tree.isBalanced() {
			t.Fatalf("delete min/max: tree is not balanced")
		}
		if !tree.is23() {
			t.Fatalf("delete min/max: tree is not a 2-3 tree")
		}
	}
}
