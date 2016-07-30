package llrb

import (
	"reflect"
	"sort"
	"testing"
)

type compInt int

func (i compInt) Compare(elem Element) int {
	if v, ok := elem.(compInt); ok {
		return int(i - v)
	}
	panic("unknown type")
}

type compInts []compInt

func (c compInts) Len() int           { return len(c) }
func (c compInts) Less(i, j int) bool { return c[i].Compare(c[j]) < 0 }
func (c compInts) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

func TestForEach(t *testing.T) {
	values := compInts{-10, -32, 100, 46, 239, 2349, 101, 0, 1}
	tree := &Tree{}
	txn := tree.Txn()
	for _, v := range values {
		txn.Insert(v)
	}
	tree = txn.Commit()

	var result compInts
	v := func(elem Element) bool {
		result = append(result, elem.(compInt))
		return false
	}

	tree.ForEach(v)
	sort.Sort(values)

	if !reflect.DeepEqual(values, result) {
		t.Fatalf("foreach: expected values %v, have %v", values, result)
	}
}

func TestRange(t *testing.T) {
	values := compInts{-10, -32, 100, 46, 239, 2349, 101, 0, 1}
	tree := &Tree{}
	txn := tree.Txn()
	for _, v := range values {
		txn.Insert(v)
	}
	tree = txn.Commit()

	var result compInts
	v := func(elem Element) bool {
		result = append(result, elem.(compInt))
		return false
	}

	sort.Sort(values)

	tree.Range(compInt(-32), compInt(2350), v)
	if !reflect.DeepEqual(values, result) {
		t.Fatalf("foreach: expected values %v, have %v", values, result)
	}
	result = result[:0]

	tree.Range(compInt(-32), compInt(2), v)
	if !reflect.DeepEqual(values[:4], result) {
		t.Fatalf("foreach: expected values %v, have %v", values[:4], result)
	}
	result = result[:0]

	tree.Range(compInt(-10), compInt(2), v)
	if !reflect.DeepEqual(values[1:4], result) {
		t.Fatalf("foreach: expected values %v, have %v", values[1:4], result)
	}
}
