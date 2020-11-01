package bst

import (
	"testing"
)

type Doc struct {
	ID int
	Evaluator
}

func (d *Doc) Value() int {
	return d.ID
}

func genTree() Tree {
	var tree Tree

	var docs = []Doc{
		{ID: 8},
		Doc{ID: 3}, Doc{ID: 10},
		Doc{ID: 1}, Doc{ID: 6}, Doc{ID: 14},
		Doc{ID: 4}, Doc{ID: 7}, Doc{ID: 13},
	}

	for _, doc := range docs {
		tree.Insert(&doc)
	}

	return tree
}

func TestTree_Insert(t *testing.T) {
	tree := genTree()

	b1 := tree.Root.L.L.Value
	b2 := tree.Root.L.R.L.Value
	b3 := tree.Root.L.R.R.Value
	b4 := tree.Root.R.R.L.Value

	var want = [4]int{1, 4, 7, 13}
	var got = [4]int{b1, b2, b3, b4}

	if got != want {
		t.Fatalf("want %v, got %v", want, got)
	}

	if tree.Len != 9 {
		t.Fatalf("want 9, got %v", tree.Len)
	}
}

func TestTree_Find(t *testing.T) {
	tree := genTree()

	_, err := tree.Find(0)
	if err == nil {
		t.Fatalf("want error, got nil")
	}

	vertex, err := tree.Find(4)
	if err != nil {
		t.Fatalf("want nil, got error")
	}

	if vertex.Value != 4 {
		t.Fatalf("want vertex = 4, got vertex = %d", vertex.Value)
	}
}
