package llrb

import (
	"fmt"
	"math/rand"
)

type Int int

func (i Int) Compare(elem Element) int {
	if v, ok := elem.(Int); ok {
		return int(i - v)
	}
	panic("unknown type")
}

func Example() {
	tree := &Tree{}
	txn := tree.Txn()

	for i := range rand.Perm(1000) {
		txn.Insert(Int(i))
	}

	tree = txn.Commit()

	elem := tree.Get(Int(500))
	fmt.Println(elem)

	// Output:
	// 500
}
