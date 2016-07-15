// Derived from https://raw.githubusercontent.com/biogo/store/master/llrb/llrb_test.go
//
// Copyright ©2012 The bíogo Authors. All rights reserved.
// Copyright ©2016 Markus Sonderegger. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"fmt"
	"reflect"
	"testing"
)

func (n *node) is23() bool {
	if n == nil {
		return true
	}

	// If the node has two children, only one of them may be red.
	// The other must be black...
	if (n.left != nil) && (n.right != nil) {
		if n.left.isRed() && n.right.isRed() {
			return false
		}
	}
	// And the red node should really should be the left one.
	if n.right.isRed() {
		return false
	}
	if n.isRed() && n.left.isRed() {
		return false
	}
	return n.left.is23() && n.right.is23()
}

func (n *node) isBalanced(black int) bool {
	if n == nil && black == 0 {
		return true
	} else if n == nil && black != 0 {
		return false
	}
	if !n.isRed() {
		black--
	}
	return n.left.isBalanced(black) && n.right.isBalanced(black)
}

func (n *node) isBST(min, max Element) bool {
	if n == nil {
		return true
	}
	if n.elem.Compare(min) < 0 || n.elem.Compare(max) > 0 {
		return false
	}
	return n.left.isBST(min, n.elem) && n.right.isBST(n.elem, max)
}

type compRune rune

func (cr compRune) Compare(r Element) int {
	return int(cr) - int(r.(compRune))
}

// makeTree builds a tree from a simplified Newick format returning the
// root node. Single letter node names only, no error checking and all
// nodes are full or leaf.
func makeTree(desc string) (n *node) {
	var build func([]rune) (*node, int)
	build = func(desc []rune) (cn *node, i int) {
		if len(desc) == 0 || desc[0] == ';' {
			return nil, 0
		}

		var c int
		cn = &node{}
		for {
			b := desc[i]
			i++
			if b == '(' {
				cn.left, c = build(desc[i:])
				i += c
				continue
			}
			if b == ',' {
				cn.right, c = build(desc[i:])
				i += c
				continue
			}
			if b == ')' {
				if cn.left == nil && cn.right == nil {
					return nil, i
				}
				continue
			}
			if b != ';' {
				cn.elem = compRune(b)
			}
			return cn, i
		}
		panic("makeTree: cannot reach")
	}

	n, _ = build([]rune(desc))
	if n.left == nil && n.right == nil {
		n = nil
	}
	return n
}

// describeTree returns a Newick format description of a tree defined
// by a node.
func describeTree(n *node, char, color bool) string {
	s := []rune(nil)

	var follow func(*node)
	follow = func(n *node) {
		children := n.left != nil || n.right != nil
		if children {
			s = append(s, '(')
		}
		if n.left != nil {
			follow(n.left)
		}
		if children {
			s = append(s, ',')
		}
		if n.right != nil {
			follow(n.right)
		}
		if children {
			s = append(s, ')')
		}
		if n.elem != nil {
			if char {
				s = append(s, rune(n.elem.(compRune)))
			} else {
				s = append(s, []rune(fmt.Sprintf("%d", n.elem))...)
			}
			//if color {
			//	s = append(s, []rune(fmt.Sprintf(" %v", n.color()))...)
			//}
		}
	}
	if n == nil {
		s = []rune("()")
	} else {
		follow(n)
	}
	s = append(s, ';')

	return string(s)
}

func TestMakeAndDescribeTree(t *testing.T) {
	desc := describeTree((*node)(nil), true, false)
	if desc != "();" {
		t.Fatalf("describe tree: expected %q, got %q", "();", desc)
	}

	for _, desc := range []string{
		"();",
		"((a,c)b,(e,g)f)d;",
	} {
		tree := makeTree(desc)
		if d := describeTree(tree, true, false); d != desc {
			t.Fatalf("make tree: expected %q, got %q", desc, d)
		}
	}

}

func TestRotateLeft(t *testing.T) {
	orig := "((a,c)b,(e,g)f)d;"
	rot := "(((a,c)b,e)d,g)f;"

	tree := makeTree(orig)
	tree = tree.rotateLeft()
	desc := describeTree(tree, true, false)
	if desc != rot {
		t.Fatalf("rotate left: expected %q, got %q", rot, desc)
	}

	rotTree := makeTree(rot)
	if !reflect.DeepEqual(rotTree, tree) {
		t.Fatalf("rotate left: original and rotated tree differ")
	}
}

func TestRotateRight(t *testing.T) {
	orig := "((a,c)b,(e,g)f)d;"
	rot := "(a,(c,(e,g)f)d)b;"

	tree := makeTree(orig)
	tree = tree.rotateRight()
	desc := describeTree(tree, true, false)
	if desc != rot {
		t.Fatalf("rotate right: expected %q, got %q", rot, desc)
	}

	rotTree := makeTree(rot)
	if !reflect.DeepEqual(rotTree, tree) {
		t.Fatalf("rotate right: original and rotated tree differ")
	}
}
